package models

type Tenant struct {
	Base
	Name             string
	Description      *string
	IsActive         bool
	IsBillingEnabled bool
}
