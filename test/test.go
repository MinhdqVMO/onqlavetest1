package test

type TestCommand struct {
	data map[string]string
}

func (cmd *TestCommand) Binding() string {
	return "TestCommand"
}

func (cmd *TestCommand) Metadata() map[string]string {
	if len(cmd.data) == 0 {
		return nil
	}
	return cmd.data
}
