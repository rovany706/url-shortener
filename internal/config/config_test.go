package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const programName = "test"

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantConfig AppConfig
	}{
		{
			"empty args",
			[]string{programName},
			AppConfig{BaseURL: defaultBaseURL, AppRunAddress: defaultAppRunAddress},
		},
		{
			"only BaseURL",
			[]string{programName, "-b", "http://test.com/"},
			AppConfig{BaseURL: "http://test.com/", AppRunAddress: defaultAppRunAddress},
		},
		{
			"only AppRunAddress",
			[]string{programName, "-a", ":8888"},
			AppConfig{BaseURL: defaultBaseURL, AppRunAddress: ":8888"},
		},
		{
			"full args",
			[]string{programName, "-a", ":8888", "-b", "http://test.com/"},
			AppConfig{BaseURL: "http://test.com/", AppRunAddress: ":8888"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualConfig, err := ParseArgs(tt.args[0], tt.args[1:])

			require.NoError(t, err)
			assert.Equal(t, &tt.wantConfig, actualConfig)
		})
	}
}
