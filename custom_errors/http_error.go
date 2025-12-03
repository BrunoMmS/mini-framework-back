package custom_errors

import "fmt"

type HttpError struct {
    Status  int
    Message string
}
func NotFound(entity string) *HttpError {
    return &HttpError{
        Status:  404,
        Message: fmt.Sprintf("%s not found", entity),
    }
}

func BadRequest(msg string) *HttpError {
    return &HttpError{
        Status:  400,
        Message: msg,
    }
}

func Unauthorized() *HttpError {
    return &HttpError{
        Status:  401,
        Message: "unauthorized",
    }
}

func (e *HttpError) Error() string {
    return fmt.Sprintf("%d %s", e.Status, e.Message)
}
