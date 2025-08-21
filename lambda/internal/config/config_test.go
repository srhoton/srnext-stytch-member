package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		errMsg  string
		check   func(t *testing.T, cfg *Config)
	}{
		{
			name: "Valid configuration with all required fields",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
				"STYTCH_ENV":        "test",
				"LOG_LEVEL":         "debug",
				"AWS_REGION":        "us-west-2",
				"REQUEST_TIMEOUT":   "60",
				"MAX_RETRIES":       "5",
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "test-project-id", cfg.StytchProjectID)
				assert.Equal(t, "test-secret", cfg.StytchSecret)
				assert.Equal(t, "test", cfg.StytchEnv)
				assert.Equal(t, "debug", cfg.LogLevel)
				assert.Equal(t, "us-west-2", cfg.AWSRegion)
				assert.Equal(t, 60, cfg.RequestTimeout)
				assert.Equal(t, 5, cfg.MaxRetries)
			},
		},
		{
			name: "Valid configuration with defaults",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "test-project-id", cfg.StytchProjectID)
				assert.Equal(t, "test-secret", cfg.StytchSecret)
				assert.Equal(t, "test", cfg.StytchEnv)
				assert.Equal(t, "info", cfg.LogLevel)
				assert.Equal(t, "us-east-1", cfg.AWSRegion)
				assert.Equal(t, 30, cfg.RequestTimeout)
				assert.Equal(t, 3, cfg.MaxRetries)
			},
		},
		{
			name: "Live environment",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
				"STYTCH_ENV":        "live",
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "live", cfg.StytchEnv)
			},
		},
		{
			name: "Missing project ID",
			envVars: map[string]string{
				"STYTCH_SECRET": "test-secret",
			},
			wantErr: true,
			errMsg:  "STYTCH_PROJECT_ID is required",
		},
		{
			name: "Missing secret",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
			},
			wantErr: true,
			errMsg:  "STYTCH_SECRET is required",
		},
		{
			name: "Invalid Stytch environment",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
				"STYTCH_ENV":        "invalid",
			},
			wantErr: true,
			errMsg:  "STYTCH_ENV must be either 'test' or 'live'",
		},
		{
			name: "Zero request timeout",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
				"REQUEST_TIMEOUT":   "0",
			},
			wantErr: true,
			errMsg:  "REQUEST_TIMEOUT must be positive",
		},
		{
			name: "Negative request timeout",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
				"REQUEST_TIMEOUT":   "-1",
			},
			wantErr: true,
			errMsg:  "REQUEST_TIMEOUT must be positive",
		},
		{
			name: "Negative max retries",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
				"MAX_RETRIES":       "-1",
			},
			wantErr: true,
			errMsg:  "MAX_RETRIES cannot be negative",
		},
		{
			name: "Invalid numeric values",
			envVars: map[string]string{
				"STYTCH_PROJECT_ID": "test-project-id",
				"STYTCH_SECRET":     "test-secret",
				"REQUEST_TIMEOUT":   "not-a-number",
				"MAX_RETRIES":       "also-not-a-number",
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				// Should use defaults for invalid numeric values
				assert.Equal(t, 30, cfg.RequestTimeout)
				assert.Equal(t, 3, cfg.MaxRetries)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearEnv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer clearEnv()

			cfg, err := Load()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
				if tt.check != nil {
					tt.check(t, cfg)
				}
			}
		})
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Returns environment value when set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "actual",
			expected:     "actual",
		},
		{
			name:         "Returns default when not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvWithDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsIntWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "Returns parsed integer when valid",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "42",
			expected:     42,
		},
		{
			name:         "Returns default when not set",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "",
			expected:     10,
		},
		{
			name:         "Returns default when invalid",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "not-a-number",
			expected:     10,
		},
		{
			name:         "Handles negative numbers",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "-5",
			expected:     -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvAsIntWithDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func clearEnv() {
	envVars := []string{
		"STYTCH_PROJECT_ID",
		"STYTCH_SECRET",
		"STYTCH_ENV",
		"LOG_LEVEL",
		"AWS_REGION",
		"REQUEST_TIMEOUT",
		"MAX_RETRIES",
	}

	for _, v := range envVars {
		os.Unsetenv(v)
	}
}
