package commands

type ParsedArgType int

const (
	Word ParsedArgType = iota
	KeyValuePair
)

type ParsedArg struct {
	Type  ParsedArgType
	Value string
	Key   string
}

type Command struct {
	Name       string
	ParsedArgs []ParsedArg
}

type CommandResponse struct {
	Result string
}
