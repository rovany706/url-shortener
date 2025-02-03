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
			*NewConfig(),
		},
		{
			"only BaseURL",
			[]string{programName, "-b", "http://test.com/"},
			*NewConfig(WithBaseURL("http://test.com/")),
		},
		{
			"only AppRunAddress",
			[]string{programName, "-a", ":8888"},
			*NewConfig(WithAppRunAddress(":8888")),
		},
		{
			"only LogLevel",
			[]string{programName, "-l", "info"},
			*NewConfig(WithLogLevel("info")),
		},
		{
			"only file storage path",
			[]string{programName, "-f", "storage.json"},
			*NewConfig(WithFileStoragePath("storage.json")),
		},
		{
			"full args",
			[]string{programName, "-a", ":8888", "-b", "http://test.com/", "-l", "debug"},
			*NewConfig(WithBaseURL("http://test.com/"), WithAppRunAddress(":8888"), WithLogLevel("debug")),
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

func TestParseArgsErr(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			"invalid BaseURL",
			[]string{programName, "-b", "http:,,test.com/"},
			ErrInvalidBaseURL,
		},
		{
			"invalid AppRunAddress",
			[]string{programName, "-a", "test:onetwothree"},
			ErrInvalidAppRunAddress,
		},
		{
			"invalid LogLevel",
			[]string{programName, "-l", "debug123"},
			ErrInvalidLogLevel,
		},
		{
			"invalid file storage path",
			[]string{programName, "-f", ""},
			ErrInvalidFileStoragePath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseArgs(tt.args[0], tt.args[1:])
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}
