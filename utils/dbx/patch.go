package dbx

import (
	"errors"
	"reflect"
)

// ApplyPatch applies the non-nil fields from a dbx struct to a target struct.
func ApplyPatch(patch any, target any) error {
	patchVal := reflect.ValueOf(patch)
	targetVal := reflect.ValueOf(target)

	if patchVal.Kind() != reflect.Ptr {
		return errors.New("patch must be a pointer")
	}
	if targetVal.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	patchVal = patchVal.Elem()
	targetVal = targetVal.Elem()

	if patchVal.Kind() != reflect.Struct {
		return errors.New("patch must be a pointer to struct")
	}
	if targetVal.Kind() != reflect.Struct {
		return errors.New("target must be a pointer to struct")
	}

	if patchVal.Type() != targetVal.Type() {
		return errors.New("patch and target must be of the same type")
	}

	patchType := patchVal.Type()

	for i := range patchVal.NumField() {
		fieldType := patchType.Field(i)
		patchField := patchVal.Field(i)

		if patchField.Kind() != reflect.Ptr {
			continue
		}

		if patchField.IsNil() {
			continue
		}

		targetField := targetVal.FieldByName(fieldType.Name)
		if !targetField.IsValid() || !targetField.CanSet() {
			continue
		}

		// Check if the target field is also a pointer type
		if targetField.Kind() == reflect.Ptr {
			// Both patch and target fields are pointers, set the entire pointer
			targetField.Set(patchField)
		} else {
			// Target field is not a pointer, set the dereferenced value
			targetField.Set(patchField.Elem())
		}

	}

	return nil
}
