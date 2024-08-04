package proto

const (
	CommandSet = "set"
	CommandGet = "get"
)

type Command interface{}

type SetCommand struct {
	Key, Value []byte
}

type GetCommand struct {
	Key []byte
}
