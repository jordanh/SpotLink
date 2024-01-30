package commands

type Command struct {
	Name string
	Args []string
}

type CommandResponse struct {
	Result string
}
