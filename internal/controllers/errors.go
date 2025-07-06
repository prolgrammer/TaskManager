package controllers

import "errors"

var (
	ErrDataBindError           = errors.New("wrong data format")
	ErrInvalidPaginationParams = errors.New("invalid pagination parameters")
)
