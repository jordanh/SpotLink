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
	Final  bool
}

type EvalContext struct {
	session      *CommandSession
	commandChan  chan Command
	responseChan chan CommandResponse
	quiet        bool
}

type EvalOption func(*EvalContext)
