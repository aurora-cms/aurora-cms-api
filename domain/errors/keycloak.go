package errors

import "errors"

var ErrKeycloakIDEmpty = errors.New("keycloak ID cannot be empty")
var ErrKeycloakIDInvalid = errors.New("keycloak ID is invalid")
