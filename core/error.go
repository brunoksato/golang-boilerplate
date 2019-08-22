package core

import (
	"reflect"

	"github.com/brunoksato/golang-boilerplate/util"
)

type DefaultError interface {
	Code() int
	Subcode() int
	Error() string
	Location() string
	Data() map[string]interface{}
	IsWarning() bool
}

const ERROR_CODE_WARNING int = 300
const ERROR_CODE_BUSINESS_ERROR int = 400
const ERROR_CODE_AUTHENTICATION_ERROR int = 401
const ERROR_CODE_PERMISSION_ERROR int = 403
const ERROR_CODE_NOT_FOUND int = 404
const ERROR_CODE_SERVER_ERROR int = 500

const ERROR_SUBCODE_UNDEFINED_IGNORE int = -1000
const ERROR_SUBCODE_UNDEFINED int = -1999

const ERROR_SUBCODE_FK int = -2000

const ERROR_SUBCODE_CREDENTIALS_INVALID int = -2001

const ERROR_SUBCODE_EMAIL int = -2002
const ERROR_SUBCODE_NAME_TAKEN int = -2003
const ERROR_SUBCODE_NAME_LENGTH int = -2004
const ERROR_SUBCODE_NAME_FORMAT int = -2005
const ERROR_SUBCODE_EMAIL_TAKEN int = -2006
const ERROR_SUBCODE_EMAIL_FORMAT int = -2007
const ERROR_SUBCODE_PASSWORD_LENGTH int = -2008
const ERROR_SUBCODE_PASSWORD_FORMAT int = -2009
const ERROR_SUBCODE_USERNAME_TAKEN int = -2010
const ERROR_SUBCODE_USERNAME_LENGTH int = -2011
const ERROR_SUBCODE_USERNAME_FORMAT int = -2012
const ERROR_SUBCODE_PHONE_TAKEN int = -2013
const ERROR_SUBCODE_PHONE_LENGTH int = -2014
const ERROR_SUBCODE_PHONE_FORMAT int = -2015

const ERROR_SUBCODE_USER_UNDERAGE int = -2800
const ERROR_SUBCODE_USER_LACKS_PERMISSION int = -2801
const ERROR_SUBCODE_OTHER_USER_LACKS_PERMISSION int = -2802

const ERROR_SUBCODE_DATABASE_UNAVAILABLE int = -2900
const ERROR_SUBCODE_SERVER_OVERLOADED int = -2910

type CoreError struct {
	ErrCode     int
	ErrSubcode  int
	Message     string
	ErrLocation string
	ErrData     map[string]interface{}
}

func (err CoreError) Code() int {
	return err.ErrCode
}

func (err CoreError) Subcode() int {
	return err.ErrSubcode
}

func (err CoreError) Error() string {
	return err.Message
}

func (err CoreError) Location() string {
	return err.ErrLocation
}

func (err CoreError) Data() map[string]interface{} {
	return err.ErrData
}

func (err CoreError) IsWarning() bool {
	return err.ErrCode == ERROR_CODE_WARNING
}

func NewDefaultError(code, subcode int, location, msg string, data map[string]interface{}) DefaultError {
	if subcode == 0 {
		if code == ERROR_CODE_BUSINESS_ERROR {
			subcode = ERROR_SUBCODE_UNDEFINED
		} else {
			subcode = ERROR_SUBCODE_UNDEFINED_IGNORE
		}
	}

	err := CoreError{
		ErrCode:     code,
		ErrSubcode:  subcode,
		Message:     msg,
		ErrLocation: location,
		ErrData:     data,
	}

	return err
}

func NewWarning(msg string, opts ...interface{}) DefaultError {
	return newElipsisError(ERROR_CODE_WARNING, msg, opts...)
}

func NewBusinessError(msg string, opts ...interface{}) DefaultError {
	return newElipsisError(ERROR_CODE_BUSINESS_ERROR, msg, opts...)
}

func NewAuthenticationError(msg string, opts ...interface{}) DefaultError {
	return newElipsisError(ERROR_CODE_AUTHENTICATION_ERROR, msg, opts...)
}

func NewPermissionError(msg string, opts ...interface{}) DefaultError {
	return newElipsisError(ERROR_CODE_PERMISSION_ERROR, msg, opts...)
}

func NewNotFoundError(msg string, opts ...interface{}) DefaultError {
	return newElipsisError(ERROR_CODE_NOT_FOUND, msg, opts...)
}

func NewServerError(msg string, opts ...interface{}) DefaultError {
	return newElipsisError(ERROR_CODE_SERVER_ERROR, msg, opts...)
}

func newElipsisError(code int, msg string, opts ...interface{}) DefaultError {
	subcode := 0
	data := map[string]interface{}{}

	for _, opt := range opts {
		switch reflect.TypeOf(opt).Kind() {
		case reflect.Int:
			subcode = opt.(int)
		case reflect.Map:
			data = opt.(map[string]interface{})
		}
	}

	return NewDefaultError(code, subcode, util.ParentCallerInfo(), msg, data)
}
