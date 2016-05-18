package models

type Err struct {
	Message string
}

func (e *Err) Error() string {
	return e.Message
}

func NewError(msg string) *Err {
	return &Err{msg}
}

func Error(err error) *Err {
	return &Err{err.Error()}
}
