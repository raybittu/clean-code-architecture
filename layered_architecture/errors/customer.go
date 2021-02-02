package errors

import "fmt"

type MissingParams struct {
	Params []string
}

func (m MissingParams) Error() string {
	return fmt.Sprintf("You have missing parameters %v", m.Params)
}

type Error string

type ErrorWithCode struct {
	ErrMsg Error
	Code   int
}

const ErrEligibility Error = "you are not eligible to be our customer"
const ErrDBQuery Error = "Wrong DB Query"
const ErrDBExec Error = "Wrong DB Execution"
const ErrInvalidID Error = "Invalid ID parameter"
const ErrDBNotFound Error = "Record Not Found"
const ErrJSONFormat Error = "Json format not correct"

func (e Error) Error() string {
	return string(e)
}

func (e ErrorWithCode) Error() string {
	return fmt.Sprintf("error is %v and code is %v", e.ErrMsg, e.Code)
}
