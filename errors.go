package fieldmask

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNilInput = errors.New("cannot apply fieldmask to nil struct")
	ErrNoStruct = errors.New("input is not a struct")
)

type errUnexpectedKind struct {
	kind reflect.Kind
}

func (e *errUnexpectedKind) Error() string {
	return fmt.Sprintf("unexpected field kind: %s", e.kind.String())
}

func IsUnexpectedKindError(err error) bool {
	var errUnexpectedKind *errUnexpectedKind
	return errors.As(err, &errUnexpectedKind)
}

type errCircularReference struct {
	reflectType reflect.Type
}

func (e *errCircularReference) Error() string {
	return fmt.Sprintf("circular reference detected for type %v", e.reflectType.Name())
}

func IsCircularReferenceError(err error) bool {
	var errCircularReference *errCircularReference
	return errors.As(err, &errCircularReference)
}

type errFieldProcessing struct {
	fieldName string
	err       error
}

func (e *errFieldProcessing) Error() string {
	return fmt.Sprintf("failed to process field %s: %v", e.fieldName, e.err)
}

func IsFieldProcessingError(err error) bool {
	var errFieldProcessing *errFieldProcessing
	return errors.As(err, &errFieldProcessing)
}
