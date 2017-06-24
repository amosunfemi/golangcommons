package utility

import (
	"fmt"
	"time"
)

//CustomErr ...
type CustomErr struct {
	ErrNo int
	When  time.Time
	Msg   string
}

//Error ...
func (err *CustomErr) Error() string {
	return fmt.Sprintf("%v [%d] %s", err.When, err.ErrNo, err.Msg)
}

//errorNumber ...
func (err CustomErr) errorNumber() int {
	return err.ErrNo
}
