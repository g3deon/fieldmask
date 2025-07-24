package fieldmask_test

import (
	"reflect"
	"testing"

	"go.g3deon.com/fieldmask"
)

func TestFieldMask_New(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  *fieldmask.FieldMask
	}{
		{
			name:  "empty paths",
			paths: []string{},
			want:  nil,
		},
		{
			name:  "non-empty paths",
			paths: []string{"field1", "field2"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1", "field2"}},
		},
		{
			name:  "duplicate paths",
			paths: []string{"field1", "field1", "field2"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1", "field2"}},
		},
		{
			name:  "empty path",
			paths: []string{"field1", ""},
			want:  &fieldmask.FieldMask{Paths: []string{"field1"}},
		},
		{
			name:  "single path",
			paths: []string{"field1"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fieldmask.New(tt.paths...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldMask_Normalize(t *testing.T) {
	tests := []struct {
		name  string
		mask  *fieldmask.FieldMask
		paths []string
		want  *fieldmask.FieldMask
	}{
		{
			name: "empty mask",
			mask: fieldmask.New(),
		},
		{
			name:  "non-empty mask",
			mask:  fieldmask.New("field1", "field2"),
			paths: []string{"field1", "field2"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1", "field2"}},
		},
		{
			name:  "duplicate paths",
			mask:  fieldmask.New("field1", "field1", "field2"),
			paths: []string{"field1", "field2"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1", "field2"}},
		},
		{
			name:  "empty path",
			mask:  fieldmask.New("field1", ""),
			paths: []string{"field1"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1"}},
		},
		{
			name:  "single path",
			mask:  fieldmask.New("field1"),
			paths: []string{"field1"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1"}},
		},
		{
			name:  "sub path",
			mask:  fieldmask.New("field1.subfield"),
			paths: []string{"field1.subfield"},
			want:  &fieldmask.FieldMask{Paths: []string{"field1.subfield"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mask.Normalize()

			if !reflect.DeepEqual(tt.mask, tt.want) {
				t.Errorf("Normalize() = %v, want %v", tt.mask, tt.want)
			}
		})
	}
}

func TestFieldMask_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		mask *fieldmask.FieldMask
		want bool
	}{
		{
			name: "empty mask",
			mask: fieldmask.New(),
			want: true,
		},
		{
			name: "non-empty mask",
			mask: fieldmask.New("field1", "field2"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mask.IsEmpty()
			if got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldMask_HasPath(t *testing.T) {
	tests := []struct {
		name string
		mask *fieldmask.FieldMask
		path string
		want bool
	}{
		{
			name: "path exists",
			mask: fieldmask.New("field1", "field2"),
			path: "field1",
			want: true,
		},
		{
			name: "path does not exist",
			mask: fieldmask.New("field1", "field2"),
			path: "field3",
			want: false,
		},
		{
			name: "empty mask",
			mask: fieldmask.New(),
			path: "field1",
			want: false,
		},
		{
			name: "sub path exists",
			mask: fieldmask.New("field1.subfield"),
			path: "field1.subfield",
			want: true,
		},
		{
			name: "sub path does not exist",
			mask: fieldmask.New("field1.subfield"),
			path: "field1.subfield.subfield",
			want: false,
		},
		{
			name: "parent path with trailing dot",
			mask: fieldmask.New("field1.subfield"),
			path: "field1",
			want: true,
		},
		{
			name: "nil mask",
			mask: nil,
			path: "field1",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mask.HasPath(tt.path)
			if got != tt.want {
				t.Errorf("HasPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldMask_HashAny(t *testing.T) {
	tests := []struct {
		name  string
		mask  *fieldmask.FieldMask
		paths []string
		want  bool
	}{
		{
			name:  "empty mask",
			mask:  fieldmask.New(),
			paths: []string{"field1", "field2"},
			want:  false,
		},
		{
			name:  "non-empty mask",
			mask:  fieldmask.New("field1", "field2"),
			paths: []string{"field1", "field2", "field3"},
			want:  true,
		},
		{
			name:  "nil mask",
			mask:  nil,
			paths: []string{"field1", "field2"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mask.HashAny(tt.paths...)
			if got != tt.want {
				t.Errorf("HashAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldMask_GetPaths(t *testing.T) {
	tests := []struct {
		name string
		mask *fieldmask.FieldMask
		want []string
	}{
		{
			name: "empty mask",
			mask: fieldmask.New(),
			want: nil,
		},
		{
			name: "non-empty mask",
			mask: fieldmask.New("field1", "field2"),
			want: []string{"field1", "field2"},
		},
		{
			name: "duplicate paths",
			mask: fieldmask.New("field1", "field1", "field2"),
			want: []string{"field1", "field2"},
		},
		{
			name: "single path",
			mask: fieldmask.New("field1"),
			want: []string{"field1"},
		},
		{
			name: "sub path",
			mask: fieldmask.New("field1.subfield"),
			want: []string{"field1.subfield"},
		},
		{
			name: "nil mask",
			mask: nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mask.GetPaths()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldMask_RemovePaths(t *testing.T) {
	tests := []struct {
		name  string
		mask  *fieldmask.FieldMask
		paths []string
		want  *fieldmask.FieldMask
	}{
		{
			name:  "remove existing paths",
			mask:  fieldmask.New("field1", "field2", "field3"),
			paths: []string{"field1", "field3"},
			want:  fieldmask.New("field2"),
		},
		{
			name:  "remove non-existing paths",
			mask:  fieldmask.New("field1", "field2"),
			paths: []string{"field3"},
			want:  fieldmask.New("field1", "field2"),
		},
		{
			name:  "remove all paths",
			mask:  fieldmask.New("field1", "field2"),
			paths: []string{"field1", "field2"},
			want:  &fieldmask.FieldMask{Paths: make([]string, 0)},
		},
		{
			name: "remove sub paths",
			mask: fieldmask.New("field1", "field1.subfield", "field2.subfield"),
			paths: []string{
				"field1",
				"field2.subfield",
			},
			want: &fieldmask.FieldMask{Paths: make([]string, 0)},
		},
		{
			name:  "remove from empty mask",
			mask:  fieldmask.New(),
			paths: []string{"field1"},
			want:  fieldmask.New(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mask.RemovePaths(tt.paths...)
			if !reflect.DeepEqual(tt.mask, tt.want) {
				t.Errorf("RemovePaths() = %v, want %v", tt.mask, tt.want)
			}
		})
	}
}

func TestFieldMask_Apply(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"`
	}

	type NoTagStruct struct {
		Field1 string
		Field2 string
	}

	type NestedStruct struct {
		Subfield string `json:"subfield"`
	}

	type ComplexStruct struct {
		Field1 NestedStruct `json:"field1"`
		Field2 string       `json:"field2"`
	}

	type PrivateFieldStruct struct {
		Public  string `json:"public"`
		private string
	}

	type PointerNestedStruct struct {
		Nested *NestedStruct `json:"nested"`
		Field  string        `json:"field"`
	}

	type MixedTagsStruct struct {
		WithTag    string `json:"with_tag"`
		WithoutTag string
		AnotherTag string `json:"another_tag"`
		NoTag      int
	}

	type NestedMixedStruct struct {
		JSONField   string `json:"json_field"`
		NoTagField  string
		Nested      MixedTagsStruct `json:"nested"`
		NoTagNested NoTagStruct
	}

	tests := []struct {
		name      string
		mask      *fieldmask.FieldMask
		input     any
		want      any
		wantError bool
	}{
		{
			name:  "apply empty mask",
			mask:  fieldmask.New(),
			input: &struct{}{},
			want:  &struct{}{},
		},
		{
			name:  "apply non-empty mask",
			mask:  fieldmask.New("field1", "field2"),
			input: &TestStruct{Field1: "value1", Field2: "value2"},
			want:  &TestStruct{Field1: "value1", Field2: "value2"},
		},
		{
			name:      "apply mask to nil input",
			mask:      fieldmask.New("field1"),
			input:     nil,
			want:      nil,
			wantError: true,
		},
		{
			name:  "apply mask with sub paths",
			mask:  fieldmask.New("field1.subfield"),
			input: &ComplexStruct{Field1: NestedStruct{Subfield: "sub-value"}, Field2: "value2"},
			want:  &ComplexStruct{Field1: NestedStruct{Subfield: "sub-value"}, Field2: ""},
		},
		{
			name:  "apply mask to non-existent field",
			mask:  fieldmask.New("nonexistent"),
			input: &TestStruct{Field1: "x", Field2: "y"},
			want:  &TestStruct{Field1: "", Field2: ""},
		},
		{
			name:  "apply mask with private field",
			mask:  fieldmask.New("private", "public"),
			input: &PrivateFieldStruct{Public: "a", private: "b"},
			want:  &PrivateFieldStruct{Public: "a", private: "b"},
		},
		{
			name:  "apply mask on struct with no json tags using field names",
			mask:  fieldmask.New("Field1", "Field2"),
			input: &NoTagStruct{Field1: "x", Field2: "y"},
			want:  &NoTagStruct{Field1: "x", Field2: "y"},
		},
		{
			name:  "apply mask on struct with no json tags using partial field names",
			mask:  fieldmask.New("Field1"),
			input: &NoTagStruct{Field1: "x", Field2: "y"},
			want:  &NoTagStruct{Field1: "x", Field2: ""},
		},
		{
			name:  "apply mask on struct with no json tags using wrong case",
			mask:  fieldmask.New("field1", "field2"),
			input: &NoTagStruct{Field1: "x", Field2: "y"},
			want:  &NoTagStruct{Field1: "", Field2: ""},
		},
		{
			name:  "apply mask on missing nested field",
			mask:  fieldmask.New("nested.subfield"),
			input: &PointerNestedStruct{Nested: nil, Field: "value"},
			want:  &PointerNestedStruct{Nested: nil, Field: ""},
		},
		{
			name:  "apply mask on struct with mixed json tags and field names",
			mask:  fieldmask.New("with_tag", "WithoutTag", "NoTag"),
			input: &MixedTagsStruct{WithTag: "tagged", WithoutTag: "notag", AnotherTag: "another", NoTag: 42},
			want:  &MixedTagsStruct{WithTag: "tagged", WithoutTag: "notag", AnotherTag: "", NoTag: 42},
		},
		{
			name:  "apply mask using field names instead of json tags on mixed struct",
			mask:  fieldmask.New("WithTag", "AnotherTag"),
			input: &MixedTagsStruct{WithTag: "tagged", WithoutTag: "notag", AnotherTag: "another", NoTag: 42},
			want:  &MixedTagsStruct{WithTag: "", WithoutTag: "", AnotherTag: "", NoTag: 0},
		},
		{
			name: "apply mask on nested struct with mixed tags",
			mask: fieldmask.New("json_field", "NoTagField", "nested.with_tag", "NoTagNested.Field1"),
			input: &NestedMixedStruct{
				JSONField:   "json_val",
				NoTagField:  "no_tag_val",
				Nested:      MixedTagsStruct{WithTag: "nested_tagged", WithoutTag: "nested_notag", AnotherTag: "nested_another", NoTag: 123},
				NoTagNested: NoTagStruct{Field1: "field1_val", Field2: "field2_val"},
			},
			want: &NestedMixedStruct{
				JSONField:   "json_val",
				NoTagField:  "no_tag_val",
				Nested:      MixedTagsStruct{WithTag: "nested_tagged", WithoutTag: "", AnotherTag: "", NoTag: 0},
				NoTagNested: NoTagStruct{Field1: "field1_val", Field2: ""},
			},
		},
		{
			name:      "handle nil interface",
			mask:      fieldmask.New("field"),
			input:     nil,
			wantError: true,
		},
		{
			name:      "handle non-pointer struct",
			mask:      fieldmask.New("field"),
			input:     TestStruct{},
			wantError: true,
		},
		{
			name: "deep nested path",
			mask: fieldmask.New("field1.subfield.deeper.even_deeper"),
			input: &ComplexStruct{
				Field1: NestedStruct{Subfield: "value"},
				Field2: "test",
			},
			want: &ComplexStruct{
				Field1: NestedStruct{Subfield: "value"},
				Field2: "",
			},
		},
		{
			name: "handle nil pointer in nested structure",
			mask: fieldmask.New("nested.subfield"),
			input: &PointerNestedStruct{
				Nested: nil,
				Field:  "value",
			},
			want: &PointerNestedStruct{
				Nested: nil,
				Field:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mask.Apply(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("Apply() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Apply() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(tt.input, tt.want) {
				t.Errorf("Apply() = %+v, want %+v", tt.input, tt.want)
			}
		})
	}
}

func TestFieldMask_Apply_Panics(t *testing.T) {
	tests := []struct {
		name        string
		mask        *fieldmask.FieldMask
		setupObj    func() interface{}
		description string
	}{
		{
			name:        "nil interface field",
			mask:        fieldmask.New("field"),
			description: "Should not panic when field is nil interface",
			setupObj: func() interface{} {
				return &struct {
					Field interface{} `json:"field"`
				}{Field: nil}
			},
		},
		{
			name:        "invalid reflection",
			mask:        fieldmask.New("field"),
			description: "Should not panic on invalid reflection target",
			setupObj: func() interface{} {
				var i interface{}
				return &i
			},
		},
		{
			name:        "empty struct pointer",
			mask:        fieldmask.New("nonexistent"),
			description: "Should not panic on empty struct with nonexistent field",
			setupObj: func() interface{} {
				return &struct{}{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Test case %q panicked: %v", tt.name, r)
				}
			}()

			input := tt.setupObj()
			tt.mask.Apply(input)
		})
	}
}
