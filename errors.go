package bizday

type bizdayError struct {
	code    string
	message string
}

func (be *bizdayError) Error() string {
	return be.code + ": " + be.message
}

func (e *bizdayError) Is(target error) bool {
	t, ok := target.(*bizdayError)
	return ok && e.code == t.code
}

var InvalidDateError = &bizdayError{
	code:    "invalid-date",
	message: "invalid date format, should be YYYY-MM-DD.",
}

var InvalidServiceError = &bizdayError{
	code:    "invalid-holiday-service",
	message: "service must not be nil",
}
