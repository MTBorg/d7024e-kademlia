package rpccommand

type RPCCommand interface {
	Execute()
	ParseOptions(options *[]string) error
}
