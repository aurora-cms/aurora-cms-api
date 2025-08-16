package mysql

import (
	"github.com/h4rdc0m/aurora-api/infrastructure/config"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		env  *config.Env
		want Config
	}{
		{
			name: "valid config",
			env: &config.Env{
				DBHost:     "localhost",
				DBPort:     "3306",
				DBUser:     "testuser",
				DBPassword: "testpassword",
				DBName:     "testdb",
			},
			want: Config{
				Host:     "localhost",
				Port:     "3306",
				Username: "testuser",
				Password: "testpassword",
				Database: "testdb",
			},
		},
		{
			name: "empty config",
			env:  &config.Env{},
			want: Config{
				Host:     "",
				Port:     "",
				Username: "",
				Password: "",
				Database: "",
			},
		},
		{
			name: "partial config",
			env: &config.Env{
				DBHost:     "127.0.0.1",
				DBPort:     "8000",
				DBUser:     "",
				DBPassword: "mypassword",
				DBName:     "",
			},
			want: Config{
				Host:     "127.0.0.1",
				Port:     "8000",
				Username: "",
				Password: "mypassword",
				Database: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig(tt.env)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
