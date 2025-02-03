package models

import "errors"

var ErrNoRecord = errors.New("models: not matching record found")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrConversionJSON = errors.New("models: couldn't convert data to json")
