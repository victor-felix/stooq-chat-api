package errors

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type ErrorCode string

type ErrorInfo struct {
	Code           ErrorCode
	Args           map[string]interface{}
	Message        string
	DefaultMessage string
}

func (e *ErrorInfo) Error() string {
	var buf bytes.Buffer
	if e.Code != "" {
		fmt.Fprintf(&buf, "<%s> ", e.Code)
	}

	if strings.TrimSpace(e.Message) == "" {
		if _, err := buf.WriteString(e.DefaultMessage); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
	} else {
		if _, err := buf.WriteString(e.Message); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
	}

	if len(e.Args) > 0 {
		if _, err := buf.WriteString(" ("); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
		if _, err := buf.WriteString(getKeyValueAsString(e.Args)); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
		if _, err := buf.WriteString(")"); err != nil {
			return fmt.Sprintf("%s : %s", string(e.Code), err.Error())
		}
	}

	return buf.String()
}

// Not Found error

type errNotFound struct {
	ErrorInfo
}

func NewErrNotFound(code ErrorCode) *errNotFound {
	err := &errNotFound{}
	err.Code = code
	err.DefaultMessage = "not found"
	return err
}

func (e *errNotFound) WithMessage(msg string) *errNotFound {
	e.Message = msg
	return e
}

func (e *errNotFound) WithMessagef(msg string, args ...interface{}) *errNotFound {
	e.Message = fmt.Sprintf(msg, args...)

	return e
}

func (e *errNotFound) WithArg(key string, value interface{}) *errNotFound {
	if e.Args == nil {
		e.Args = make(map[string]interface{})
	}
	e.Args[key] = value
	return e
}

func (e *errNotFound) Error() string {
	return e.ErrorInfo.Error()
}

func (e *errNotFound) ErrNotFound() {}

func (e *errNotFound) GetMessage() string {
	return e.Message
}

func (e *errNotFound) GetCode() ErrorCode {
	return e.Code
}

func (e *errNotFound) GetArgs() map[string]interface{} {
	return e.Args
}

// Unauthorized error

type errUnauthorized struct {
	ErrorInfo
}

func NewErrUnauthorized(code ErrorCode) *errUnauthorized {
	err := &errUnauthorized{}
	err.Code = code
	err.DefaultMessage = "unauthorized"
	return err
}

func (e *errUnauthorized) WithMessage(msg string) *errUnauthorized {
	e.Message = msg
	return e
}

func (e *errUnauthorized) WithMessagef(msg string, args ...interface{}) *errUnauthorized {
	e.Message = fmt.Sprintf(msg, args...)

	return e
}

func (e *errUnauthorized) WithArg(key string, value interface{}) *errUnauthorized {
	if e.Args == nil {
		e.Args = make(map[string]interface{})
	}
	e.Args[key] = value
	return e
}

func (e *errUnauthorized) Error() string {
	return e.ErrorInfo.Error()
}

func (e *errUnauthorized) ErrUnauthorized() {}

func (e *errUnauthorized) GetMessage() string {
	return e.Message
}

func (e *errUnauthorized) GetCode() ErrorCode {
	return e.Code
}

func (e *errUnauthorized) GetArgs() map[string]interface{} {
	return e.ErrorInfo.Args
}

// Rule Not Satisfied error

type errRuleNotSatisfied struct {
	ErrorInfo
}

func NewErrRuleNotSatisfied(code ErrorCode) *errRuleNotSatisfied {
	err := &errRuleNotSatisfied{}
	err.Code = code
	err.DefaultMessage = "internal rule was not satisfied"
	return err
}

func (e *errRuleNotSatisfied) WithMessage(msg string) *errRuleNotSatisfied {
	e.Message = msg
	return e
}

func (e *errRuleNotSatisfied) WithMessagef(msg string, args ...interface{}) *errRuleNotSatisfied {
	e.Message = fmt.Sprintf(msg, args...)

	return e
}

func (e *errRuleNotSatisfied) WithArg(key string, value interface{}) *errRuleNotSatisfied {
	if e.Args == nil {
		e.Args = make(map[string]interface{})
	}
	e.Args[key] = value
	return e
}

func (e *errRuleNotSatisfied) Error() string {
	return e.ErrorInfo.Error()
}

func (e *errRuleNotSatisfied) ErrRulesNotSatisfied() {}

func (e *errRuleNotSatisfied) GetMessage() string {
	return e.Message
}

func (e *errRuleNotSatisfied) GetCode() ErrorCode {
	return e.Code
}

func (e *errRuleNotSatisfied) GetArgs() map[string]interface{} {
	return e.ErrorInfo.Args
}

// Invalid Argument error

type errInvalidArgument struct {
	ErrorInfo
}

func NewErrInvalidArgument(code ErrorCode) *errInvalidArgument {
	err := &errInvalidArgument{}
	err.Code = code
	err.DefaultMessage = "invalid argument"
	return err
}

func (e *errInvalidArgument) WithMessage(msg string) *errInvalidArgument {
	e.Message = msg
	return e
}

func (e *errInvalidArgument) WithArg(key string, value interface{}) *errInvalidArgument {
	if e.Args == nil {
		e.Args = make(map[string]interface{})
	}
	e.Args[key] = value
	return e
}

func (e *errInvalidArgument) Error() string {
	return e.ErrorInfo.Error()
}

func (e *errInvalidArgument) ErrInvalidArgument() {}

func (e *errInvalidArgument) GetMessage() string {
	return e.Message
}

func (e *errInvalidArgument) GetCode() ErrorCode {
	return e.Code
}

func (e *errInvalidArgument) GetArgs() map[string]interface{} {
	return e.Args
}

func getKeyValueAsString(args map[string]interface{}) string {
	values := make([]string, 0, len(args))

	for key, value := range args {
		values = append(values, fmt.Sprintf("%s: %v", key, value))
	}

	sort.Strings(values)
	return strings.Join(values, ", ")
}
