package marathonw

// Error is a trivial error object
type Error struct {
	Message string `json:"message"`
}

// Error returns the error description
func (e Error) Error() string {
	return e.Message
}

var (
	ErrServiceNameCollision = &Error{Message: "a service name collision has been dectected: %s"}
	ErrServiceLabelNotFound = &Error{Message: "label not found in marathon application"}
)
