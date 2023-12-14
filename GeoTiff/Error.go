package GeoTiff

import (
	"errors"
	"fmt"
)

type GeoError struct {
	Err      error
	Function string
	Msg      string
}

func (ge GeoError) Error() string {
	return fmt.Sprintf("[%s]: %v\n", ge.Function, ge.Err)
}

type GeoErrorOptions func(ge *GeoError)

func WithFunction(Function string) GeoErrorOptions {
	return func(ge *GeoError) {
		ge.Function = Function
	}
}
func WithMsg(Msg string) GeoErrorOptions {
	return func(ge *GeoError) {
		ge.Msg = Msg
	}
}
func WithError(err error) GeoErrorOptions {
	return func(ge *GeoError) {
		ge.Err = err
	}
}
func WithErrorText(text string) GeoErrorOptions {
	return func(ge *GeoError) {
		ge.Err = errors.New(text)
	}
}

func NewGeoErrorCreator(Function string) func(opts ...GeoErrorOptions) GeoError {
	return func(opts ...GeoErrorOptions) GeoError {
		ge := GeoError{
			Err:      nil,
			Function: Function,
			Msg:      "",
		}
		for _, opt := range opts {
			opt(&ge)
		}
		return ge
	}
}
