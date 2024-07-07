package mdbxsql

type Model interface {
	PrimaryKey() interface{}
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}
