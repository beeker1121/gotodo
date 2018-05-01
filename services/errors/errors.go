package errors

// NewParamError returns a new ParamError.
func NewParamError(name string, err error) *ParamError {
	return &ParamError{
		Name:      name,
		ErrorType: err,
	}
}

// ParamError defines an error with a parameter passed to a service method.
type ParamError struct {
	Name      string
	ErrorType error
}

// Error implements the error interface.
func (pe *ParamError) Error() string {
	return pe.ErrorType.Error()
}

// NewParamErrors returns a new ParamErrors.
//
// The main goal of this is to have the service methods themselves capable
// of returning multiple errors, in this case, parameter errors. The API
// specifically can then just iterate through these errors, placing them
// into an errors array for JSON output.
func NewParamErrors(pe ...*ParamError) *ParamErrors {
	pes := &ParamErrors{}

	for _, v := range pe {
		pes.Add(v)
	}

	return pes
}

// ParamErrors defines a slice of parameter errors.
type ParamErrors []*ParamError

// Add appends a ParamError.
func (pes *ParamErrors) Add(pe *ParamError) {
	*pes = append(*pes, pe)
}

// Length returns the number of parameter errors.
func (pes *ParamErrors) Length() int {
	return len(*pes)
}

// Error implements the error interface.
func (pes *ParamErrors) Error() string {
	var b []byte
	for _, pe := range *pes {
		b = append(b, pe.Error()...)
		b = append(b, '\n')
	}
	return string(b)
}
