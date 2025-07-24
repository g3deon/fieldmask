package fieldmask_test

import (
	"testing"

	"go.g3deon.com/fieldmask"
)

func BenchmarkFieldmask_Apply(b *testing.B) {
	fm := &fieldmask.FieldMask{Paths: []string{"field1", "field2", "field3"}}
	type testStruct struct {
		Field1 string
		Field2 int
		Field3 bool
		Field4 float64
	}
	s := &testStruct{}

	b.ReportAllocs()

	for b.Loop() {
		*s = testStruct{}
		fm.Apply(s)
	}
}
