package v4

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
)

// Bounds resolution maps a validate tag onto the domain of the flag's type,
// so `min`/`max` compare in the unit the author meant:
//
//   - string kinds (incl. Enum/LogLevel/LogFormat): length bounds
//     (existing behavior; min/max mean characters).
//   - duration types (Duration, time.Duration): duration bounds
//     ("min=10s,max=1m" compares durations, never nanoseconds against "10").
//   - textual parsed types (URL, Email, HostPort, FilePath): length bounds;
//     their Stringer text is not a number, so the numeric validators would
//     fail with a wrong-error class on every run.
//   - everything else (numeric kinds, Port, ...): numeric bounds.
func validateExprForType(tag FlagTag) string {
	if tag.Validate == "" {
		return ""
	}

	switch {
	case isDurationType(tag.Type):
		return durationBoundsRules(tag.Validate)
	case tag.Type == nil || tag.Type.Kind() == reflect.String || isTextualParsedType(tag.Type):
		return stringLengthRules(tag.Validate)
	default:
		return tag.Validate
	}
}

func isDurationType(t reflect.Type) bool {
	return t == reflect.TypeFor[Duration]() || t == reflect.TypeFor[time.Duration]()
}

func isTextualParsedType(t reflect.Type) bool {
	switch t {
	case reflect.TypeFor[URL](),
		reflect.TypeFor[Email](),
		reflect.TypeFor[HostPort](),
		reflect.TypeFor[FilePath]():
		return true
	default:
		return false
	}
}

// durationBoundsRules rewrites numeric bounds to their duration equivalents
// for a duration-typed field, so validate:"min=10s" is a duration bound.
func durationBoundsRules(expr string) string {
	return rewriteNumericBounds(expr, "dmin", "dmax")
}

// rewriteNumericBounds renames the min/max rules in a validate tag expression,
// leaving every other rule untouched.
func rewriteNumericBounds(expr, minName, maxName string) string {
	parts := strings.Split(expr, ",")

	for i, part := range parts {
		name, param, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}

		switch strings.TrimSpace(name) {
		case "min":
			parts[i] = minName + "=" + param
		case "max":
			parts[i] = maxName + "=" + param
		}
	}

	return strings.Join(parts, ",")
}

// Construction-time default validation: an impossible default must fail
// NewFlagRegistry — where the mistake is made — not the first run.

// validateDefaults checks every tag default against that tag's values and
// validate rules. A validate tag naming a validator that is not resolvable at
// construction time (instance-scoped validators are often registered after
// NewFlagRegistry) is skipped here; run-time validation still covers it.
func (r *FlagRegistry) validateDefaults() error {
	for _, tag := range r.tags {
		if tag.Default == "" {
			continue
		}

		if err := validateTagDefault(tag, r.validators); err != nil {
			return err
		}
	}

	return nil
}

func validateTagDefault(tag FlagTag, vr *validatorRegistry) error {
	if len(tag.Values) > 0 && !slices.Contains(tag.Values, tag.Default) {
		return fmt.Errorf(
			"%w: flag %q default %q, allowed: %v: %w",
			ErrInvalidFlagDefault,
			tag.Name,
			tag.Default,
			tag.Values,
			ErrInvalidEnum,
		)
	}

	if tag.Validate == "" || !rulesResolvable(tag.Validate, vr) {
		return nil
	}

	rules, err := parseValidateRulesWithRegistry(validateExprForType(tag), vr)
	if err != nil {
		return fmt.Errorf("%w: flag %q: %w", ErrInvalidFlagDefault, tag.Name, err)
	}

	for _, rule := range rules {
		if err := rule.Validate(tag.Default); err != nil {
			return fmt.Errorf(
				"%w: flag %q default %q: %w",
				ErrInvalidFlagDefault,
				tag.Name,
				tag.Default,
				err,
			)
		}
	}

	return nil
}

// rulesResolvable reports whether every validator named in the tag exists in
// the instance or global registry. Unknown names are legal at construction
// (late instance registration), so an unresolvable name disables the
// construction-time check for that tag instead of failing it.
func rulesResolvable(tag string, vr *validatorRegistry) bool {
	for part := range strings.SplitSeq(tag, ",") {
		name, _, _ := strings.Cut(strings.TrimSpace(part), "=")
		if name == "" {
			continue
		}

		if vr != nil {
			if _, ok := vr.lookup(name); ok {
				continue
			}
		}

		if _, ok := lookupValidator(name); !ok {
			return false
		}
	}

	return true
}
