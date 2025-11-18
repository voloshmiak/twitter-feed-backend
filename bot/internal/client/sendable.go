package client

type Sendable interface {
	MarshalJSON() ([]byte, error)
}
