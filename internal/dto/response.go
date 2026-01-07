package dto

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(code int, status, message string, body interface{}) *Response {
	return &Response{
		Code:    code,
		Status:  status,
		Message: message,
		Body:    body,
	}
}

func Error(code int, status, message string, err interface{}) *Response {
	return &Response{
		Code:    code,
		Status:  status,
		Message: message,
		Error:   err,
	}
}
