package errors

import "errors"

var ErrPageKeyEmpty = errors.New("page key cannot be empty")
var ErrPageKeyTooLong = errors.New("page key cannot be longer than 255 characters")
var ErrPageKeyInvalid = errors.New("page key can only contain alphanumeric characters, underscores, and hyphens")
var ErrPageTypeInvalid = errors.New("page type is invalid")
var ErrPageVersionEmpty = errors.New("page version cannot be empty")
var ErrPageVersionTitleEmpty = errors.New("page version title cannot be empty")
var ErrPageVersionInvalidVersion = errors.New("page version is invalid")
var ErrPageBlockEmpty = errors.New("page block cannot be empty")
var ErrInvalidContentType = errors.New("content type is invalid")
var ErrInvalidBlockKey = errors.New("block key is invalid")
var ErrInvalidPageType = errors.New("invalid page type")
var ErrInvalidPageVersionModel = errors.New("invalid page version model")
