package secure

import "fmt"

type SafeError struct {
	Code     string
	UserMsg  string
	Internal error
	Metadata map[string]string
}

func (e *SafeError) Error() string {
	return e.UserMsg
}

func (e *SafeError) LogString() string {
	return fmt.Sprintf("Code: %s | Msg: %s | Cause: %v | Meta: %v",
		e.Code, e.UserMsg, e.Internal, e.Metadata)
}

func NewAuthFailed(public string, internal error, safeMeta map[string]string) *SafeError {
	return &SafeError{
		Code:     "AUTH_FAILED",
		UserMsg:  public,
		Internal: internal,
		Metadata: safeMeta,
	}
}
