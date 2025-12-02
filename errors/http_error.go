package errors

import "fmt"

type HttpError struct {
    Status  int
    Message string
}

func (e *HttpError) Error() string {
    return fmt.Sprintf("%d %s", e.Status, e.Message)
}
