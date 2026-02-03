package helper

import "errors"

var (
	ErrParsingFailed = errors.New("parsing failed")
	ErrInvalidDate = errors.New("invalid date format")
	ErrAttrNotFound = errors.New("attribute not found")
	ErrTimeOut = errors.New("parsing timeout after 30 seconds")
	ErrBadRequest = errors.New("couldn't get html")
)