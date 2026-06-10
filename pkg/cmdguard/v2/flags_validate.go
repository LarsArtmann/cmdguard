package v2

import (
	"fmt"
	"maps"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// FlagValidator validates a flag value and returns an error if invalid.
// The value parameter is the string representation of the flag value.
type FlagValidator func(value string) error

// validatorRegistry holds named validators that can be referenced via validate tags.
type validatorRegistry struct {
	mu         sync.RWMutex
	validators map[string]FlagValidator
}

// globalValidators is the package-level validator registry.
var globalValidators = newValidatorRegistry()

func newValidatorRegistry() *validatorRegistry {
	registry := &validatorRegistry{
		validators: make(map[string]FlagValidator),
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

	r.validators[name] = validator
}

func (r *validatorRegistry) clone() *validatorRegistry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c := &validatorRegistry{
		validators: make(map[string]FlagValidator, len(r.validators)),
	}

	maps.Copy(c.validators, r.validators)

	return c
}

func (r *validatorRegistry) lookup(name string) (FlagValidator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.validators[name]

	return v, ok
}

// RegisterValidator adds a named validator to the global defaults template.
// New FlagRegistries will include this validator. For per-instance registration,
// use FlagRegistry.RegisterFlagValidator.
func RegisterValidator(name string, validator FlagValidator) {
	globalValidators.register(name, validator)
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

	_, err := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("%w: %q is not a valid email", ErrInvalidEmail, value)
	}

	return nil
}

func validateURL(value string) error {
	if value == "" {
		return nil
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("%w: %q is not a valid URL", ErrInvalidURL, value)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%w: %q must include scheme and host", ErrInvalidURL, value)
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

	if len(val) < minLen {
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

	if len(val) > maxLen {
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
		return fmt.Errorf("%w: value must not be empty", ErrValueEmpty)
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
		return fmt.Errorf("parsing validation rules: %w", err)
	}

	for _, rule := range rules {
		if err := rule.Validate(value); err != nil {
			return fmt.Errorf("validation rule %q failed: %w", rule.Name, err)
		}
	}

	return nil
}
