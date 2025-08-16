package entities

import (
	"errors"
	"time"
)

// TemplateSettingID represents a template setting identifier
type TemplateSettingID struct {
	value uint
}

// NewTemplateSettingID creates a new TemplateSettingID
func NewTemplateSettingID(id uint) TemplateSettingID {
	return TemplateSettingID{value: id}
}

// Value returns the ID value
func (t TemplateSettingID) Value() uint {
	return t.value
}

// TemplateSetting represents a template setting entity
type TemplateSetting struct {
	id           TemplateSettingID
	templateID   TemplateID
	settingKey   string
	settingValue string
	canOverride  bool
	createdAt    time.Time
	updatedAt    time.Time
}

// NewTemplateSetting creates a new TemplateSetting entity
func NewTemplateSetting(templateID TemplateID, settingKey, settingValue string, canOverride bool) (*TemplateSetting, error) {
	if settingKey == "" {
		return nil, errors.New("setting key is required")
	}

	if settingValue == "" {
		return nil, errors.New("setting value is required")
	}

	now := time.Now()

	return &TemplateSetting{
		templateID:   templateID,
		settingKey:   settingKey,
		settingValue: settingValue,
		canOverride:  canOverride,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ID returns the template setting ID
func (t *TemplateSetting) ID() TemplateSettingID {
	return t.id
}

// TemplateID returns the template ID
func (t *TemplateSetting) TemplateID() TemplateID {
	return t.templateID
}

// SettingKey returns the setting key
func (t *TemplateSetting) SettingKey() string {
	return t.settingKey
}

// SettingValue returns the setting value
func (t *TemplateSetting) SettingValue() string {
	return t.settingValue
}

// CanOverride returns whether the setting can be overridden
func (t *TemplateSetting) CanOverride() bool {
	return t.canOverride
}

// CreatedAt returns the creation time
func (t *TemplateSetting) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the last update time
func (t *TemplateSetting) UpdatedAt() time.Time {
	return t.updatedAt
}

// UpdateSettingValue updates the setting value
func (t *TemplateSetting) UpdateSettingValue(value string) error {
	if value == "" {
		return errors.New("setting value cannot be empty")
	}

	t.settingValue = value
	return nil
}

// SetCanOverride sets whether the setting can be overridden
func (t *TemplateSetting) SetCanOverride(canOverride bool) {
	t.canOverride = canOverride
}

// SetID sets the template setting ID (used by repository when loading from database)
func (t *TemplateSetting) SetID(id TemplateSettingID) {
	t.id = id
}

// SetTimestamps sets the timestamps (used by repository when loading from database)
func (t *TemplateSetting) SetTimestamps(createdAt, updatedAt time.Time) {
	t.createdAt = createdAt
	t.updatedAt = updatedAt
}
