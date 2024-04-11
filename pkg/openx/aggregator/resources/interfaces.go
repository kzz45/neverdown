package resources

type Object interface {
	Marshal() (dAtA []byte, err error)
	Unmarshal(dAtA []byte) error
}
