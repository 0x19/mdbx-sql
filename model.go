package mdbxsql

type Model interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}
