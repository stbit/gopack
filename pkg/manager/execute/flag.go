package execute

type CommandsFlag []string

func (c *CommandsFlag) String() string {
	return "exec commands"
}

func (c *CommandsFlag) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func (c *CommandsFlag) Len() int {
	return len(*c)
}
