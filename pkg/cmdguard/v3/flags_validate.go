package v3

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

// share returns a copy-on-write view of this registry.
// The returned instance reads from this registry's map until the first write,
// at which point it clones lazily. This avoids the clone cost for the common
// case where no per-instance customization is used.
func (r *validatorRegistry) share() *validatorRegistry {
	root := r
	if !r.owned && r.parent != nil {
		root = r.parent
	}

	root.mu.RLock()
	defer root.mu.RUnlock()

	return &validatorRegistry{
		owned:  false,
		parent: root,
	}
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
	minStr, val, ok := strings.Cut(value, ":")
	if !ok {
		return fmt.Errorf("%w: minlen requires format \"min:value\"", ErrInvalidValidatorParam)
	}

	minLen, err := strconv.Atoi(minStr)
	if err != nil {
		return fmt.Errorf("%w: minlen: invalid integer %q", ErrInvalidValidatorParam, minStr)
	}

	if utf8.RuneCountInString(val) < minLen {
		return fmt.Errorf("%w: %q must be at least %d characters", ErrValueTooShort, val, minLen)
	}

	return nil
}

func validateMaxLen(value string) error {
	maxStr, val, ok := strings.Cut(value, ":")
	if !ok {
		return fmt.Errorf("%w: maxlen requires format \"max:value\"", ErrInvalidValidatorParam)
	}

	maxLen, err := strconv.Atoi(maxStr)
	if err != nil {
		return fmt.Errorf("%w: maxlen: invalid integer %q", ErrInvalidValidatorParam, maxStr)
	}

	if utf8.RuneCountInString(val) > maxLen {
		return fmt.Errorf("%w: %q must be at most %d characters", ErrValueTooLong, val, maxLen)
	}

	return nil
}

func validateMin(value string) error {
	minStr, val, ok := strings.Cut(value, ":")
	if !ok {
		return fmt.Errorf("%w: min requires format \"min:value\"", ErrInvalidValidatorParam)
	}

	minVal, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		return fmt.Errorf("%w: min: invalid number %q", ErrInvalidValidatorParam, minStr)
	}

	actual, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fmt.Errorf("%w: min: value %q is not a number", ErrInvalidValidatorParam, val)
	}

	if actual < minVal {
		return fmt.Errorf("%w: %v must be at least %v", ErrValueTooSmall, actual, minVal)
	}

	return nil
}

func validateMax(value string) error {
	maxStr, val, ok := strings.Cut(value, ":")
	if !ok {
		return fmt.Errorf("%w: max requires format \"max:value\"", ErrInvalidValidatorParam)
	}

	maxVal, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		return fmt.Errorf("%w: max: invalid number %q", ErrInvalidValidatorParam, maxStr)
	}

	actual, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fmt.Errorf("%w: max: value %q is not a number", ErrInvalidValidatorParam, val)
	}

	if actual > maxVal {
		return fmt.Errorf("%w: %v must be at most %v", ErrValueTooLarge, actual, maxVal)
	}

	return nil
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
