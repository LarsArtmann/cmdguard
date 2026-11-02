package v4

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unicode"
)

// Config files, flags, and struct tags are three spellings of the same setting,
// and cmdguard has to accept all of them without guessing. This file is the
// single definition of that vocabulary:
//
//   - A config-file key addresses a struct field by its json tag when present,
//     otherwise by the Go field name.
//   - Any spelling that folds to the same value is accepted, so "log-level",
//     "log_level", "LOG-LEVEL", and "LogLevel" all address the same field. That
//     matches what encoding/json/v2's case-insensitive matching already did.
//   - A flag name is an accepted alias, so the value a user types after
//     --log-level is the value they can write in the file.
//   - Keys nothing accepts are reported, not ignored.

// configKeyEntry is one accepted key spelling and the struct value it reaches.
type configKeyEntry struct {
	// canonical is the spelling written back into the JSON before unmarshalling.
	canonical string
	node      *configKeyNode
}

// configKeyObject is the accepted key vocabulary of one JSON object: the root
// config struct, a nested config group, or an array element.
type configKeyObject struct {
	exact     map[string]*configKeyEntry
	folded    map[string]*configKeyEntry
	spellings []string
}

// configKeyNode is the vocabulary of one config value. Exactly one of object,
// element, and open is meaningful; the zero value describes a leaf (a scalar,
// an opaque value type, or a field cmdguard does not look inside).
type configKeyNode struct {
	object  *configKeyObject
	element *configKeyNode
	open    bool
}

// configUnknownKey is a config-file key that no field of the config struct
// accepts, together with the closest accepted spelling at the same level.
type configUnknownKey struct {
	Path       string
	Suggestion string
}

// configKeyName returns the canonical config-file key for a struct field: the
// json tag name when set, otherwise the Go field name.
func configKeyName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return field.Name
	}

	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return field.Name
	}

	return name
}

// foldConfigKey normalizes a key for spelling-insensitive lookup: case is
// ignored and separators are dropped. "log-level", "log_level", and "LogLevel"
// all fold to "loglevel".
func foldConfigKey(key string) string {
	var b strings.Builder

	b.Grow(len(key))

	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			continue
		}

		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}

// add registers one accepted spelling for entry. Conflicting writes to the same
// folded key are disarmed rather than resolved by declaration order: a config
// file that is ambiguous for the struct is a bug worth surfacing, not guessing
// about.
func (o *configKeyObject) add(spelling string, entry *configKeyEntry) {
	if _, exists := o.exact[spelling]; exists {
		return
	}

	o.exact[spelling] = entry

	folded := foldConfigKey(spelling)

	if existing, ok := o.folded[folded]; ok && existing != entry {
		o.folded[folded] = nil

		return
	}

	o.folded[folded] = entry
}

// merge copies every spelling of other into o, preserving o's own definitions.
func (o *configKeyObject) merge(other *configKeyObject) {
	for _, spelling := range other.spellings {
		o.add(spelling, other.exact[spelling])
	}

	o.spellings = append(o.spellings, other.spellings...)
}

// lookup resolves a config-file key to its entry. Exact spellings win so an
// explicit json tag is never shadowed by a folded alias.
func (o *configKeyObject) lookup(key string) (*configKeyEntry, bool) {
	if entry, ok := o.exact[key]; ok {
		return entry, true
	}

	entry, ok := o.folded[foldConfigKey(key)]
	if !ok || entry == nil {
		return nil, false
	}

	return entry, true
}

// addEmbedded promotes an anonymous struct's fields into o, mirroring how
// encoding/json flattens embedded structs. It reports whether field was an
// embedded config group.
func (o *configKeyObject) addEmbedded(field reflect.StructField) bool {
	if !field.Anonymous {
		return false
	}

	t := field.Type
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct || !shouldRecurseInto(t) {
		return false
	}

	o.merge(buildConfigKeyObject(t))

	return true
}

// buildConfigKeyObject derives the accepted vocabulary of a struct type.
func buildConfigKeyObject(t reflect.Type) *configKeyObject {
	obj := &configKeyObject{
		exact:  make(map[string]*configKeyEntry, t.NumField()),
		folded: make(map[string]*configKeyEntry, t.NumField()),
	}

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() || field.Tag.Get("json") == "-" {
			continue
		}

		if obj.addEmbedded(field) {
			continue
		}

		canonical := configKeyName(field)
		entry := &configKeyEntry{canonical: canonical, node: buildConfigKeyNode(field.Type)}

		obj.add(canonical, entry)

		if flag := field.Tag.Get("flag"); flag != "" && flag != "-" {
			obj.add(flag, entry)
		}

		obj.spellings = append(obj.spellings, canonical)
	}

	return obj
}

// buildConfigKeyNode derives the vocabulary of a value type. Structs that
// cmdguard treats as opaque values (durations, enums, ports, ...) and scalars
// are leaves.
func buildConfigKeyNode(t reflect.Type) *configKeyNode {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() { //nolint:exhaustive // every other kind is a leaf
	case reflect.Struct:
		if !shouldRecurseInto(t) {
			return &configKeyNode{}
		}

		return &configKeyNode{object: buildConfigKeyObject(t)}
	case reflect.Slice, reflect.Array:
		elem := buildConfigKeyNode(t.Elem())
		if elem.object == nil && elem.element == nil && !elem.open {
			return &configKeyNode{}
		}

		return &configKeyNode{element: elem}
	case reflect.Map:
		return &configKeyNode{open: true}
	default:
		return &configKeyNode{}
	}
}

// configKeyTree returns the root vocabulary node for a config value.
func configKeyTree(cfg any) (*configKeyNode, error) {
	v, err := derefPointerToStruct(cfg)
	if err != nil {
		return nil, err
	}

	return buildConfigKeyNode(v.Type()), nil
}

// joinConfigPath builds the dotted path shown to a user when a nested key is
// wrong. Array elements carry an index: "servers[0].port".
func joinConfigPath(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + "." + key
}

// suggestConfigKey returns the accepted spelling closest to key, comparing the
// folded forms so separators and case do not hide an obvious typo.
func suggestConfigKey(spellings []string, key string) string {
	foldedKey := foldConfigKey(key)

	best := ""
	bestDist := maxEditDistance + 1

	for _, spelling := range spellings {
		dist := editDistance(foldedKey, foldConfigKey(spelling))
		if dist < bestDist {
			bestDist = dist
			best = spelling
		}
	}

	if bestDist > maxEditDistance {
		return ""
	}

	return best
}

// scanConfigObject rewrites every recognized key to its canonical spelling and
// records the keys that nothing accepts. Unknown keys are dropped, which is
// exactly what json.Unmarshal did with them, so tolerant mode keeps its
// behavior while strict mode gains an accurate list.
func scanConfigObject(
	obj map[string]jsontext.Value,
	node *configKeyObject,
	prefix string,
	unknown *[]configUnknownKey,
) map[string]jsontext.Value {
	out := make(map[string]jsontext.Value, len(obj))

	for key, value := range obj {
		entry, ok := node.lookup(key)
		if !ok {
			*unknown = append(*unknown, configUnknownKey{
				Path:       joinConfigPath(prefix, key),
				Suggestion: suggestConfigKey(node.spellings, key),
			})

			continue
		}

		path := joinConfigPath(prefix, entry.canonical)
		out[entry.canonical] = scanConfigValue(value, entry.node, path, unknown)
	}

	return out
}

// scanConfigValue descends into value when the config type expects an object or
// an array of objects, and returns the value untouched otherwise.
func scanConfigValue(
	value jsontext.Value,
	node *configKeyNode,
	path string,
	unknown *[]configUnknownKey,
) jsontext.Value {
	if node == nil || node.open {
		return value
	}

	if node.object != nil {
		var obj map[string]jsontext.Value
		if json.Unmarshal(value, &obj) == nil {
			return encodeConfigObject(scanConfigObject(obj, node.object, path, unknown))
		}
	}

	if node.element != nil {
		var arr []jsontext.Value
		if json.Unmarshal(value, &arr) == nil {
			out := make([]jsontext.Value, 0, len(arr))

			for i, elem := range arr {
				elemPath := fmt.Sprintf("%s[%d]", path, i)
				out = append(out, scanConfigValue(elem, node.element, elemPath, unknown))
			}

			return encodeConfigArray(out)
		}
	}

	return value
}

// encodeConfigObject renders a normalized object with sorted keys so the
// normalized bytes are deterministic.
func encodeConfigObject(obj map[string]jsontext.Value) jsontext.Value {
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	var b strings.Builder

	b.WriteByte('{')

	for i, key := range keys {
		if i > 0 {
			b.WriteByte(',')
		}

		quoted, err := jsontext.AppendQuote(nil, key)
		if err != nil {
			// Only invalid UTF-8 in a config key can reach here, and the
			// replacement-character form is still valid JSON.
			quoted, _ = json.Marshal(key) //nolint:errchkjson // string marshal cannot fail
		}

		b.Write(quoted)
		b.WriteByte(':')
		b.Write(obj[key])
	}

	b.WriteByte('}')

	return jsontext.Value(b.String())
}

// encodeConfigArray renders a normalized array.
func encodeConfigArray(values []jsontext.Value) jsontext.Value {
	var b strings.Builder

	b.WriteByte('[')

	for i, value := range values {
		if i > 0 {
			b.WriteByte(',')
		}

		b.Write(value)
	}

	b.WriteByte(']')

	return jsontext.Value(b.String())
}

// normalizeConfigKeys rewrites the config file's keys to canonical spellings
// and reports the ones no field accepts.
func normalizeConfigKeys(
	raw map[string]jsontext.Value,
	cfg any,
) (jsontext.Value, []configUnknownKey, error) {
	tree, err := configKeyTree(cfg)
	if err != nil {
		return nil, nil, err
	}

	if tree.object == nil {
		return nil, nil, nil //nolint:nilnil // non-struct config is rejected by ParseFlagTags
	}

	var unknown []configUnknownKey

	out := scanConfigObject(raw, tree.object, "", &unknown)

	return encodeConfigObject(out), unknown, nil
}

// unknownConfigKeyError reports every unrecognized key, each with the nearest
// accepted spelling when one is close enough to be the likely typo.
func unknownConfigKeyError(unknown []configUnknownKey) error {
	slices.SortFunc(unknown, func(a, b configUnknownKey) int {
		return strings.Compare(a.Path, b.Path)
	})

	errs := make([]error, 0, len(unknown))

	for _, key := range unknown {
		if key.Suggestion != "" {
			errs = append(
				errs,
				fmt.Errorf("%w %q (did you mean %q?)", ErrConfigUnknownKey, key.Path, key.Suggestion),
			)

			continue
		}

		errs = append(errs, fmt.Errorf("%w %q", ErrConfigUnknownKey, key.Path))
	}

	return errors.Join(errs...)
}
