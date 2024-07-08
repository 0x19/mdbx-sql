package mdbxsql

type Model interface {
	PrimaryKey() []byte
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}
