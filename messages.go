package tgmux

type Messages struct {
	NoCommand     string `json:"noCommand"`
	InternalError string `json:"internalError"`
}

var defaultMessages = &Messages{
	NoCommand:     "No such command",
	InternalError: "Internal error, try again",
}
