package value_objects

type NullableString struct {
	value *string
}

func NewNullableString(value string) NullableString {
	if value == "" {
		return NullableString{}
	}
	return NullableString{value: &value}
}

func (ns NullableString) Value() *string {
	return ns.value
}

func (ns NullableString) IsEmpty() bool {
	return ns.value == nil
}

func (ns NullableString) String() string {
	if ns.IsEmpty() {
		return ""
	}
	return *ns.value
}

func (ns NullableString) Equals(other NullableString) bool {
	return ns.Value() == other.Value()
}
