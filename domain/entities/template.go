package entities

import (
	"errors"
	"time"
)

type TemplateID struct {
	value uint64
}

// NewTemplateID creates a new TemplateID instance with the specified unsigned integer value.
func NewTemplateID(id uint64) TemplateID {
	return TemplateID{value: id}
}

// Value retrieves the internal `value` field of the TemplateID.
func (t TemplateID) Value() uint64 {
	return t.value
}

type Template struct {
	id          TemplateID
	name        string
	description *string
	filePath    string
	enabled     bool
	createdAt   time.Time
	updatedAt   time.Time
	settings    []*TemplateSetting
}

func NewTemplate(name, filePath string, description *string) (*Template, error) {
	if name == "" {
		return nil, errors.New("template name is required")
	}

	if filePath == "" {
		return nil, errors.New("template file path is required")
	}

	now := time.Now()

	return &Template{
		name:        name,
		description: description,
		filePath:    filePath,
		enabled:     true,
		createdAt:   now,
		updatedAt:   now,
		settings:    make([]*TemplateSetting, 0),
	}, nil
}

// ID returns the template ID
func (t *Template) ID() TemplateID {
	return t.id
}

// Name returns the template name
func (t *Template) Name() string {
	return t.name
}

// Description returns the template description
func (t *Template) Description() *string {
	return t.description
}

// FilePath returns the template file path
func (t *Template) FilePath() string {
	return t.filePath
}

// IsEnabled returns whether the template is enabled
func (t *Template) IsEnabled() bool {
	return t.enabled
}

// CreatedAt returns the creation time
func (t *Template) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the last update time
func (t *Template) UpdatedAt() time.Time {
	return t.updatedAt
}

// Settings returns the template settings
func (t *Template) Settings() []*TemplateSetting {
	return t.settings
}

// UpdateName updates the template name
func (t *Template) UpdateName(name string) error {
	if name == "" {
		return errors.New("template name cannot be empty")
	}

	t.name = name

	return nil
}

// UpdateDescription updates the template description
func (t *Template) UpdateDescription(description *string) {
	t.description = description
}

// UpdateFilePath updates the template file path
func (t *Template) UpdateFilePath(filePath string) error {
	if filePath == "" {
		return errors.New("template file path cannot be empty")
	}

	t.filePath = filePath
	return nil
}

// Enable enables the template
func (t *Template) Enable() {
	t.enabled = true
}

// Disable disables the template
func (t *Template) Disable() {
	t.enabled = false
}

// AddSetting adds a template setting
func (t *Template) AddSetting(setting *TemplateSetting) error {
	if setting == nil {
		return errors.New("setting cannot be nil")
	}

	// Check if setting with same key already exists
	for _, existingSetting := range t.settings {
		if existingSetting.SettingKey() == setting.SettingKey() {
			return errors.New("setting with this key already exists")
		}
	}

	t.settings = append(t.settings, setting)
	return nil
}

// RemoveSetting removes a template setting
func (t *Template) RemoveSetting(settingID TemplateSettingID) {
	for i, setting := range t.settings {
		if setting.ID().Value() == settingID.Value() {
			t.settings = append(t.settings[:i], t.settings[i+1:]...)
			break
		}
	}
}

// SetID sets the template ID (used by repository when loading from database)
func (t *Template) SetID(id TemplateID) {
	t.id = id
}

// SetTimestamps sets the timestamps (used by repository when loading from database)
func (t *Template) SetTimestamps(createdAt, updatedAt time.Time) {
	t.createdAt = createdAt
	t.updatedAt = updatedAt
}
