package fieldmask

import (
	"reflect"
	"testing"
)

func TestTypeDescriptor_Apply(t *testing.T) {
	type Nested struct {
		SubField string `json:"subfield"`
	}
	type TestStruct struct {
		Field1     string `json:"field1"`
		Field2     int    `json:"field2"`
		Nested     Nested `json:"nested"`
		Unexported int
	}

	tests := []struct {
		name     string
		input    TestStruct
		paths    []string
		expected TestStruct
	}{
		{
			name:  "keep single field",
			input: TestStruct{Field1: "value", Field2: 42, Nested: Nested{SubField: "nestedValue"}},
			paths: []string{"field1"},
			expected: TestStruct{
				Field1: "value",
			},
		},
		{
			name:  "clear unselected fields",
			input: TestStruct{Field1: "value", Field2: 42},
			paths: []string{"field2"},
			expected: TestStruct{
				Field2: 42,
			},
		},
		{
			name:  "nested structure",
			input: TestStruct{Field1: "value", Nested: Nested{SubField: "nestedValue"}},
			paths: []string{"nested.subfield"},
			expected: TestStruct{
				Nested: Nested{SubField: "nestedValue"},
			},
		},
		{
			name:     "empty paths should clear all",
			input:    TestStruct{Field1: "value", Field2: 42},
			paths:    []string{},
			expected: TestStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := getTypeDescriptor(reflect.TypeOf(tt.input))
			if err != nil {
				t.Fatalf("failed to get descriptor: %v", err)
			}
			value := reflect.ValueOf(&tt.input).Elem()
			if err := desc.apply(value, tt.paths, map[uintptr]bool{}); err != nil {
				t.Fatalf("apply failed: %v", err)
			}
			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("apply result mismatch. got %v, want %v", tt.input, tt.expected)
			}
		})
	}
}

func TestTypeDescriptor_GetTypeDescriptor(t *testing.T) {
	type TestStruct struct {
		Field string
	}
	type Recursive struct {
		Child *Recursive
	}

	tests := []struct {
		name        string
		input       reflect.Type
		expectError bool
	}{
		{
			name:  "simple structure",
			input: reflect.TypeOf(TestStruct{}),
		},
		{
			name:        "circular reference",
			input:       reflect.TypeOf(Recursive{}),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getTypeDescriptor(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("getTypeDescriptor() error = %v, wantError = %v", err, tt.expectError)
			}
		})
	}
}

func TestTypeDescriptor_BuildDescriptor(t *testing.T) {
	type Nested struct {
		SubField string `json:"subfield"`
	}
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
		Nested Nested `json:"nested"`
		Ignore string `json:"-"`
	}

	tests := []struct {
		name        string
		input       reflect.Type
		expectError bool
		expected    []string
	}{
		{
			name:     "parse struct with fields",
			input:    reflect.TypeOf(TestStruct{}),
			expected: []string{"field1", "field2", "nested"},
		},
		{
			name:        "no exported fields",
			input:       reflect.TypeOf(struct{ unexported int }{}),
			expectError: false,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, err := buildDescriptor(tt.input, map[reflect.Type]bool{})
			if (err != nil) != tt.expectError {
				t.Errorf("buildDescriptor() error = %v, wantError %v", err, tt.expectError)
			}
			if desc != nil {
				var gotTags []string
				for tag := range desc.fields {
					gotTags = append(gotTags, tag)
				}
				if !reflect.DeepEqual(gotTags, tt.expected) {
					t.Errorf("fields mismatch, got %v, want %v", gotTags, tt.expected)
				}
			}
		})
	}
}

func TestTypeDescriptor_BuildPathMaps(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		expectedKeep  map[string]struct{}
		expectedPaths map[string][]string
	}{
		{
			name:         "simple paths",
			input:        []string{"field1", "field2.subfield"},
			expectedKeep: map[string]struct{}{"field1": {}},
			expectedPaths: map[string][]string{
				"field2": {"subfield"},
			},
		},
		{
			name:          "empty paths",
			input:         nil,
			expectedKeep:  map[string]struct{}{},
			expectedPaths: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepMap, nestedPaths := buildPathMaps(tt.input)
			if !reflect.DeepEqual(keepMap, tt.expectedKeep) {
				t.Errorf("keepMap mismatch. got %v, want %v", keepMap, tt.expectedKeep)
			}
			if !reflect.DeepEqual(nestedPaths, tt.expectedPaths) {
				t.Errorf("nestedPaths mismatch. got %v, want %v", nestedPaths, tt.expectedPaths)
			}
		})
	}
}

func TestTypeDescriptor_ParseJSONTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty tag",
			input:    "",
			expected: "",
		},
		{
			name:     "simple tag",
			input:    "field",
			expected: "field",
		},
		{
			name:     "tag with omitempty",
			input:    "field,omitempty",
			expected: "field",
		},
		{
			name:     "tag with multiple options",
			input:    "field,omitempty,string",
			expected: "field",
		},
		{
			name:     "tag with empty options",
			input:    "field,",
			expected: "field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseJSONTag(tt.input)
			if got != tt.expected {
				t.Errorf("parseJSONTag(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTypeDescriptor_GetZero(t *testing.T) {
	tests := []struct {
		name     string
		input    reflect.Type
		expected interface{}
	}{
		{
			name:     "int type",
			input:    reflect.TypeOf(0),
			expected: 0,
		},
		{
			name:     "string type",
			input:    reflect.TypeOf(""),
			expected: "",
		},
		{
			name:     "bool type",
			input:    reflect.TypeOf(true),
			expected: false,
		},
		{
			name:     "slice type",
			input:    reflect.TypeOf([]int{}),
			expected: []int(nil),
		},
		{
			name:     "map type",
			input:    reflect.TypeOf(map[string]int{}),
			expected: map[string]int(nil),
		},
		{
			name:     "struct type",
			input:    reflect.TypeOf(struct{ Name string }{}),
			expected: struct{ Name string }{},
		},
		{
			name:     "pointer type",
			input:    reflect.TypeOf((*int)(nil)),
			expected: (*int)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getZero(tt.input).Interface()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getZero(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
