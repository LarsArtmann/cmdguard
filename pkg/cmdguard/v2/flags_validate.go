package v2

import (
	"fmt"
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

// RegisterValidator adds a named validator to the global registry.
// Safe for concurrent use via internal sync.RWMutex.
// Must be called before CLI execution to ensure validators are available.
func RegisterValidator(name string, validator FlagValidator) {
	globalValidators.mu.Lock()
	defer globalValidators.mu.Unlock()

	globalValidators.validators[name] = validator
}

// lookupValidator finds a validator by name.
func lookupValidator(name string) (FlagValidator, bool) {
	globalValidators.mu.RLock()
	defer globalValidators.mu.RUnlock()

	v, ok := globalValidators.validators[name]

	return v, ok
}

// runValidateTag runs all validators specified in a validate tag.
// The tag format is a comma-separated list: "email,min=5,max=100".
// Returns the first validation error encountered.
func runValidateTag(tag, value string) error {
	rules, err := parseValidateRules(tag)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if err := rule.Validate(value); err != nil {
			return err
		}
	}

	return nil
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

// parseValidateRulesWithRegistry parses a validate tag using instance validators first.
// If instance is nil or the validator is not found, falls back to the global registry.
func parseValidateRulesWithRegistry(tag string, instance *validatorRegistry) ([]validateRule, error) {
	var rules []validateRule

	for part := range strings.SplitSeq(tag, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		name, param, _ := strings.Cut(part, "=")

		var validator FlagValidator
		var ok bool

		// Instance-scoped first
		if instance != nil {
			instance.mu.RLock()
			validator, ok = instance.validators[name]
			instance.mu.RUnlock()
		}

		// Fallback to global
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

func validateRegex(value string) error {
	pattern, val, ok := strings.Cut(value, ":")
	if !ok {
		return fmt.Errorf("%w: regex requires format \"pattern:value\"", ErrInvalidValidatorParam)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("%w: regex: invalid pattern %q", ErrInvalidValidatorParam, pattern)
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
func validateFieldByKind(field reflect.Value, tag FlagTag) error {
	if tag.Validate == "" {
		return nil
	}

	strValue := formatFieldValue(field)

	return runValidateTag(tag.Validate, strValue)
}

// formatFieldValue converts a reflect.Value to its string representation for validation.
func formatFieldValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(field.Bool())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", field.Complex())
	case reflect.Array, reflect.Slice:
		return fmt.Sprintf("%v", field.Interface())
	case reflect.Map, reflect.Struct:
		return fmt.Sprintf("%v", field.Interface())
	case reflect.Pointer, reflect.Interface:
		if field.Elem().IsValid() {
			return formatFieldValue(field.Elem())
		}

		return ""
	case reflect.Chan, reflect.Func:
		return fmt.Sprintf("%v", field.Interface())
	case reflect.Invalid:
		return ""
	case reflect.Uintptr, reflect.UnsafePointer:
		return fmt.Sprintf("%v", field.Interface())
	default:
		return ""
	}
}
