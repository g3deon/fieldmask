package fieldmask

import (
	"reflect"
	"strings"
	"sync"
)

const (
	jsonTagSeparator = ","
	jsonTagIgnore    = "-"
)

var (
	descriptorCache sync.Map
	zeroCache       sync.Map
)

type (
	typeDescriptor struct {
		fields map[string]*fieldDescriptor
	}

	fieldDescriptor struct {
		tag   string
		index []int
		child *typeDescriptor
	}
)

// apply updates the struct fields based on the provided paths, zeroing out fields not specified in the path list.
// It uses the visited map to handle circular references and avoids processing unaddressable values.
// Returns an error if any issue arises during recursive field processing.
func (d *typeDescriptor) apply(value reflect.Value, paths []string, visited map[uintptr]bool) error {
	if !value.CanAddr() {
		return nil
	}

	addr := value.UnsafeAddr()
	if visited[addr] {
		return nil
	}
	visited[addr] = true

	keepMap, nestedPaths := buildPathMaps(paths)
	for tag, desc := range d.fields {
		fieldValue := value.FieldByIndex(desc.index)
		if !fieldValue.CanSet() {
			continue
		}

		_, keep := keepMap[tag]
		if desc.child != nil {
			if sub, ok := nestedPaths[tag]; ok {
				if err := desc.child.apply(fieldValue, sub, visited); err != nil {
					return err
				}
				continue
			}
		}
		if !keep {
			fieldValue.Set(getZero(fieldValue.Type()))
		}
	}

	return nil
}

// getTypeDescriptor retrieves or builds a typeDescriptor for a given reflect.Type, caching the result for future use.
// It dereferences pointer types to their underlying element type and handles circular references during descriptor creation.
// Returns the cached or newly built typeDescriptor, or an error if descriptor creation fails.
func getTypeDescriptor(t reflect.Type) (*typeDescriptor, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if cached, ok := descriptorCache.Load(t); ok {
		return cached.(*typeDescriptor), nil
	}

	desc, err := buildDescriptor(t, map[reflect.Type]bool{})
	if err != nil {
		return nil, err
	}

	descriptorCache.Store(t, desc)
	return desc, nil
}

// buildDescriptor constructs a typeDescriptor for the given reflect.Type, including details for its exported fields.
// It skips unexported fields and fields with a JSON tag set to "-".
// Recursive calls are made for nested struct types, with circular references handled via the visited map.
// Returns the constructed typeDescriptor or an error if a circular reference is detected.
func buildDescriptor(t reflect.Type, visited map[reflect.Type]bool) (*typeDescriptor, error) {
	if visited[t] {
		return nil, &errCircularReference{reflectType: t}
	}

	visited[t] = true

	desc := &typeDescriptor{fields: make(map[string]*fieldDescriptor, t.NumField())}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == jsonTagIgnore {
			continue
		}

		tagName := parseJSONTag(jsonTag)
		if tagName == "" {
			tagName = field.Name
		}

		fd := &fieldDescriptor{
			tag:   tagName,
			index: field.Index,
		}

		if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) {
			ft := field.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			child, err := buildDescriptor(ft, visited)
			if err != nil {
				return nil, err
			}
			fd.child = child
		}

		desc.fields[tagName] = fd
	}
	return desc, nil
}

// buildPathMaps processes the given paths to separate top-level keys and nested paths into respective maps.
// It returns a map of top-level keys (`keepMap`) and a map of parent-to-child paths (`nestedPaths`).
// Each key in `keepMap` represents a top-level path, and `nestedPaths` organizes sub-paths by their parent keys.
func buildPathMaps(paths []string) (map[string]struct{}, map[string][]string) {
	keepMap := make(map[string]struct{}, len(paths))
	nestedPaths := make(map[string][]string)

	for _, p := range paths {
		parts := strings.SplitN(p, pathSeparator, 2)
		if len(parts) == 1 {
			keepMap[parts[0]] = struct{}{}
		} else {
			nestedPaths[parts[0]] = append(nestedPaths[parts[0]], parts[1])
		}
	}
	return keepMap, nestedPaths
}

// parseJSONTag extracts the field name from a JSON struct tag, ignoring additional options like "omitempty".
func parseJSONTag(tag string) string {
	if idx := strings.Index(tag, jsonTagSeparator); idx != -1 {
		return tag[:idx]
	}
	return tag
}

// getZero retrieves a cached zero value for a given type or creates and stores it if not present.
// This optimization reduces allocations caused by repeated reflect.Zero calls.
func getZero(t reflect.Type) reflect.Value {
	if v, ok := zeroCache.Load(t); ok {
		return v.(reflect.Value)
	}
	z := reflect.Zero(t)
	zeroCache.Store(t, z)
	return z
}
