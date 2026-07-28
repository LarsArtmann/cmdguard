package v4

import (
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

// FlagValidator validates a flag value and returns an error if invalid.
// The value parameter is the string representation of the flag value.
type FlagValidator func(value string) error

// validatorRegistry holds named validators that can be referenced via validate tags.
type validatorRegistry struct {
	mu         sync.RWMutex
	validators map[string]FlagValidator
	owned      bool               // true if this instance owns its map (can mutate)
	parent     *validatorRegistry // nil when owned; shared source for COW reads
}

func (r *validatorRegistry) isCowOwned() bool {
	return r.owned
}

func (r *validatorRegistry) cowParent() (*validatorRegistry, bool) {
	return r.parent, r.parent != nil
}

func (r *validatorRegistry) cowLock() {
	r.mu.RLock()
}

func (r *validatorRegistry) cowUnlock() {
	r.mu.RUnlock()
}

// globalValidators is the package-level validator registry.
var globalValidators = newValidatorRegistry()

func newValidatorRegistry() *validatorRegistry {
	registry := &validatorRegistry{
		validators: make(map[string]FlagValidator),
		owned:      true,
	}

	registry.registerBuiltins()

	return registry
}

func (r *validatorRegistry) registerBuiltins() {
	r.validators["email"] = validateEmail
	r.validators["url"] = validateURL
	r.validators["minlen"] = validateMinLen
	r.validators["maxlen"] = validateMaxLen
	r.validators["min"] = validateMin
	r.validators["max"] = validateMax
	r.validators["regex"] = validateRegex
	r.validators["nonempty"] = validateNonEmpty
}

func (r *validatorRegistry) register(name string, validator FlagValidator) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.owned && r.parent != nil {
		r.parent.mu.RLock()
		r.validators = maps.Clone(r.parent.validators)
		r.parent.mu.RUnlock()
		r.parent = nil
		r.owned = true
	}

	r.validators[name] = validator
}

// cowNode is the subset of a copy-on-write registry needed to walk to its root.
// Both typeRegistry and validatorRegistry satisfy this interface; see cowRoot.
type cowNode[T any] interface {
	isCowOwned() bool
	cowParent() (T, bool)
	cowLock()
	cowUnlock()
}

// cowRoot walks from r to its owning root: returns the parent if r is a non-owned
// view and a parent exists, otherwise r itself. Caller must hold the returned
// root's read lock for the duration of any read.
func cowRoot[T cowNode[T]](r T) T {
	if !r.isCowOwned() {
		if p, ok := r.cowParent(); ok {
			return p
		}
	}

	return r
}

// cowShare creates a copy-on-write view of r, holding r's root read lock while
// the constructor builds the new non-owned instance from the root. The lock is
// held so that the parent pointer and root fields captured by construct remain
// stable across the share.
func cowShare[T cowNode[T]](r T, construct func(root T) T) T {
	root := cowRoot(r)

	root.cowLock()
	defer root.cowUnlock()

	return construct(root)
}

// share returns a copy-on-write view of this registry.
// The returned instance reads from this registry's map until the first write,
// at which point it clones lazily. This avoids the clone cost for the common
// case where no per-instance customization is used.
func (r *validatorRegistry) share() *validatorRegistry {
	return cowShare(r, func(root *validatorRegistry) *validatorRegistry {
		return &validatorRegistry{
			owned:  false,
			parent: root,
		}
	})
}

func (r *validatorRegistry) lookup(name string) (FlagValidator, bool) {
	r.mu.RLock()

	if r.owned {
		defer r.mu.RUnlock()

		v, ok := r.validators[name]

		return v, ok
	}

	parent := r.parent
	r.mu.RUnlock()

	return parent.lookup(name)
}

// RegisterValidator adds a named validator to the global defaults template.
// New FlagRegistries will include this validator. For per-instance registration,
// use FlagRegistry.RegisterFlagValidator.
// Returns an error if name is empty or validator is nil.
func RegisterValidator(name string, validator FlagValidator) error {
	if name == "" {
		return fmt.Errorf("%w: validator name is empty", ErrServiceRegistration)
	}

	if validator == nil {
		return fmt.Errorf("%w: validator is nil for name %q", ErrServiceRegistration, name)
	}

	globalValidators.register(name, validator)

	return nil
}

// lookupValidator finds a validator by name in the global registry.
func lookupValidator(name string) (FlagValidator, bool) {
	return globalValidators.lookup(name)
}

// validateRule is a parsed validator rule with its parameter.
type validateRule struct {
	Name      string
	Parameter string
	Validate  FlagValidator
}

// parseValidateRules parses a validate tag into individual rules.
// Format: "email,min=5,max=100".
// Returns an error if an unknown validator name is encountered.
func parseValidateRules(tag string) ([]validateRule, error) {
	return parseValidateRulesWithRegistry(tag, nil)
}

// parseValidateRulesWithRegistry parses a validate tag using the instance registry.
func parseValidateRulesWithRegistry(
	tag string,
	instance *validatorRegistry,
) ([]validateRule, error) {
	var rules []validateRule

	for part := range strings.SplitSeq(tag, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		name, param, _ := strings.Cut(part, "=")

		var (
			validator FlagValidator
			ok        bool
		)

		if instance != nil {
			validator, ok = instance.lookup(name)
		}

		if !ok {
			validator, ok = lookupValidator(name)
		}

		if !ok {
			return nil, fmt.Errorf("%w: %q", ErrUnknownValidator, name)
		}

		if param != "" {
			paramValidator := validator
			validator = func(value string) error {
				return paramValidator(param + ":" + value)
			}
		}

		rules = append(rules, validateRule{
			Name:      name,
			Parameter: param,
			Validate:  validator,
		})
	}

	return rules, nil
}

// Built-in validators

func validateEmail(value string) error {
	if value == "" {
		return nil
	}

	_, err := ParseEmail(value)
	if err != nil {
		return fmt.Errorf("%w: %q is not a valid email", ErrInvalidEmail, value)
	}

	return nil
}

func validateURL(value string) error {
	if value == "" {
		return nil
	}

	_, err := ParseURL(value)
	if err != nil {
		return fmt.Errorf("%w: %q is not a valid URL", ErrInvalidURL, value)
	}

	return nil
}

func validateMinLen(value string) error {
	minLen, val, err := parseLabeledValue(value, "minlen")
	if err != nil {
		return err
	}

	if utf8.RuneCountInString(val) < minLen {
		return fmt.Errorf("%w: %q must be at least %d characters", ErrValueTooShort, val, minLen)
	}

	return nil
}

func validateMaxLen(value string) error {
	maxLen, val, err := parseLabeledValue(value, "maxlen")
	if err != nil {
		return err
	}

	if utf8.RuneCountInString(val) > maxLen {
		return fmt.Errorf("%w: %q must be at most %d characters", ErrValueTooLong, val, maxLen)
	}

	return nil
}

// parseLabeledValue splits a "param:value" string into an integer parameter and the value.
// Returns ErrInvalidValidatorParam-wrapped errors tagged with label.
func parseLabeledValue(value, label string) (int, string, error) {
	paramStr, val, ok := strings.Cut(value, ":")
	if !ok {
		return 0, "", fmt.Errorf("%w: %s requires format \"%s:value\"", ErrInvalidValidatorParam, label, label)
	}

	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, "", fmt.Errorf("%w: %s: invalid integer %q", ErrInvalidValidatorParam, label, paramStr)
	}

	return param, val, nil
}

func validateMin(value string) error {
	minVal, actual, err := parseLabeledFloat(value, "min")
	if err != nil {
		return err
	}

	if actual < minVal {
		return fmt.Errorf("%w: %v must be at least %v", ErrValueTooSmall, actual, minVal)
	}

	return nil
}

func validateMax(value string) error {
	maxVal, actual, err := parseLabeledFloat(value, "max")
	if err != nil {
		return err
	}

	if actual > maxVal {
		return fmt.Errorf("%w: %v must be at most %v", ErrValueTooLarge, actual, maxVal)
	}

	return nil
}

// parseLabeledFloat splits a "param:value" string into a float64 parameter and a float64 value.
// Returns ErrInvalidValidatorParam-wrapped errors tagged with label.
func parseLabeledFloat(value, label string) (float64, float64, error) {
	paramStr, val, ok := strings.Cut(value, ":")
	if !ok {
		return 0, 0, fmt.Errorf("%w: %s requires format \"%s:value\"", ErrInvalidValidatorParam, label, label)
	}

	param, err := strconv.ParseFloat(paramStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %s: invalid number %q", ErrInvalidValidatorParam, label, paramStr)
	}

	actual, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %s: value %q is not a number", ErrInvalidValidatorParam, label, val)
	}

	return param, actual, nil
}

// regexCache caches compiled regex patterns to avoid recompilation on every validate call.
// The cache is unbounded; in practice the set of patterns is finite and developer-defined
// (from struct tags), so it does not grow indefinitely. If user-derived patterns are ever
// supported, this should be replaced with a bounded LRU cache.
var regexCache sync.Map // map[string]*regexp.Regexp

func validateRegex(value string) error {
	pattern, val, ok := strings.Cut(value, ":")
	if !ok {
		return fmt.Errorf("%w: regex requires format \"pattern:value\"", ErrInvalidValidatorParam)
	}

	cached, loaded := regexCache.Load(pattern)

	var re *regexp.Regexp

	if loaded {
		re = cached.(*regexp.Regexp) //nolint:forcetypeassert // only stores *regexp.Regexp
	} else {
		var err error

		re, err = regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("%w: regex: invalid pattern %q", ErrInvalidValidatorParam, pattern)
		}

		regexCache.Store(pattern, re)
	}

	if !re.MatchString(val) {
		return fmt.Errorf("%w: %q does not match pattern %q", ErrValuePatternMismatch, val, pattern)
	}

	return nil
}

func validateNonEmpty(value string) error {
	if value == "" {
		return fmt.Errorf("value=%q: %w: value must not be empty", value, ErrValueEmpty)
	}

	return nil
}

// validateFieldByKind runs type-appropriate validation for a reflected field value.
func validateFieldByKind(field reflect.Value, tag FlagTag, vr *validatorRegistry) error {
	if tag.Validate == "" {
		return nil
	}

	strValue := formatFieldValue(field)

	return runValidateTagWithRegistry(tag.Validate, strValue, vr)
}

// runValidateTagWithRegistry runs all validators specified in a validate tag using the given registry.
func runValidateTagWithRegistry(tag, value string, vr *validatorRegistry) error {
	rules, err := parseValidateRulesWithRegistry(tag, vr)
	if err != nil {
		return fmt.Errorf("tag=%q, value=%q: parsing validation rules: %w", tag, value, err)
	}

	for _, rule := range rules {
		err := rule.Validate(value)
		if err != nil {
			return fmt.Errorf("value=%q: validation rule %q failed: %w", value, rule.Name, err)
		}
	}

	return nil
}
