package app

type command struct {
	name    string
	handler func() error
}

var commandList = []command{
	{"encode", cmdEncode},
	{"decode", cmdDecode},
}
