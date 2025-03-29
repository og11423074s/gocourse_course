package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("name is required")
var ErrStartDateRequired = errors.New("startDate is required")
var ErrEndDateRequired = errors.New("endDate is required")

type ErrorNotFound struct {
	CourseID string
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("course %s doesn't exist", e.CourseID)
}
