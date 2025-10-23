package abi

type Encode interface {
	Encode() ([]byte, error)
}

type Decode interface {
	Decode([]byte) error
}

type EmptyTuple struct{}

func (e EmptyTuple) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (e *EmptyTuple) Decode(data []byte) error {
	return nil
}
