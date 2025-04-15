package internals

type ErrorNotFound struct {
	Err error
}

func (e ErrorNotFound) Error() string {
	if e.Err == nil {
		return "Not found"
	}

	return e.Err.Error()
}
