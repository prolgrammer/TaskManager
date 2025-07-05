package usecases

import "errors"

var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityAlreadyExists = errors.New("entity already exists")
