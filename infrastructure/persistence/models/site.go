package models

type Site struct {
	Base
	Name             string
	Description      *string
	Domain           string
	TitleTemplate    *string
	Enabled          bool
	TemplateID       uint64
	Template         Template
	TenantID         uint64
	Tenant           Tenant
	SettingOverrides []TemplateSettingOverride
	Pages            []Page
}
