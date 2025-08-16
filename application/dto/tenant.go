package dto

type CreateTenantRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
}

type UpdateTenantRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
}
