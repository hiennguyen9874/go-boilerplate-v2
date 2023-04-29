package newPointer

func NewBoolean(value bool) *bool {
	b := value
	return &b
}
