package fieldmask

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// FieldMask enables selective field updates by specifying dot-notation paths.
type FieldMask struct {
	Paths []string `json:"paths"`
}

func (f *FieldMask) String() string {
	if f.IsEmpty() {
		return "FieldMask{}"
	}

	return fmt.Sprintf("FieldMask{Paths: %v}", strings.Join(f.Paths, ", "))
}

// IsEmpty reports whether the FieldMask is empty.
func (f *FieldMask) IsEmpty() bool {
	return f == nil || len(f.Paths) == 0
}

// Normalize removes empty and duplicate paths from the FieldMask, ensuring the Paths slice contains only unique, valid paths.
func (f *FieldMask) Normalize() {
	if f.IsEmpty() {
		return
	}

	normalized := removeEmptyPaths(f.Paths)
	normalized = removeDuplicatePaths(normalized)
	f.Paths = normalized
}

// HasPath checks if the specified path is in the FieldMask.
func (f *FieldMask) HasPath(path string) bool {
	if f.IsEmpty() {
		return false
	}

	prefix := path + "."
	for _, p := range f.Paths {
		if p == path || strings.HasPrefix(p, prefix) {
			return true
		}
	}

	return false
}

// HashAny checks if any of the given paths are present in the FieldMask and returns true if at least one match is found.
func (f *FieldMask) HashAny(paths ...string) bool {
	if f.IsEmpty() {
		return false
	}

	for _, p := range paths {
		if f.HasPath(p) {
			return true
		}
	}

	return false
}

// GetPaths returns a copy of the paths to prevent external mutation.
func (f *FieldMask) GetPaths() []string {
	if f.IsEmpty() {
		return nil
	}

	return slices.Clone(f.Paths)
}

// RemovePaths deletes the specified paths from the FieldMask.
func (f *FieldMask) RemovePaths(paths ...string) {
	if f.IsEmpty() {
		return
	}

	for _, p := range paths {
		newPaths := make([]string, 0, len(f.Paths))
		prefix := p + "."
		for _, existing := range f.Paths {
			if existing == p || strings.HasPrefix(existing, prefix) {
				continue
			}
			newPaths = append(newPaths, existing)
		}
		f.Paths = newPaths
	}
}

// Apply zeros to all struct fields except those specified in f.Paths.
func (f *FieldMask) Apply(i any) error {
	if f.IsEmpty() {
		return nil
	}

	if i == nil {
		return ErrNilInput
	}

	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNilInput
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrNoStruct
	}

	td, err := getTypeDescriptor(v.Type())
	if err != nil {
		return err
	}

	return td.apply(v, f.Paths, make(map[uintptr]bool))
}

// New creates a FieldMask with the given paths. Returns nil for empty input.
func New(paths ...string) *FieldMask {
	if len(paths) == 0 {
		return nil
	}

	fm := &FieldMask{
		Paths: paths,
	}
	fm.Normalize()
	return fm
}
