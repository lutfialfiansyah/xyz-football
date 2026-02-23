package apperror

// AppError is a structured error that carries bilingual messages and an HTTP status code.
type AppError struct {
	Code int
	En   string
	Id   string
}

func (e *AppError) Error() string {
	return e.En
}

// New creates a new AppError with HTTP status code and bilingual messages.
func New(code int, en, id string) *AppError {
	return &AppError{
		Code: code,
		En:   en,
		Id:   id,
	}
}
