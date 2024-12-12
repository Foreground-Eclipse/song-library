package response

type Response struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(data interface{}) Response {
	return Response{
		Status: StatusOK,
		Data:   data,
	}
}

func Error(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
