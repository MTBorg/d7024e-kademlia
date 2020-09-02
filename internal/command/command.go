package command

type Command interface {
	Execute() (string, error)

	// Parse the options (i.e. words after command) and set related fields in
	// the struct
	ParseOptions(options []string) error

	PrintUsage()
}
