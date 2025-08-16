package models

type Template struct {
	Base
	Name        string
	Description *string
	FilePath    string
	Enabled     bool
}

type TemplateSetting struct {
	Base
	TemplateID   string
	SettingKey   string
	SettingValue string
	CanOverride  bool
}

type TemplateSettingOverride struct {
	Base
	SiteID            string
	TemplateSettingID string
	SettingValue      string
}
