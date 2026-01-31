package response

type Body struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Body {
	return Body{
		Status: StatusOK,
	}
}

func Error(msg string) Body {
	return Body{
		Status: StatusError,
		Error:  msg,
	}
}
