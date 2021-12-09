package config_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jamieaitken/cgs/config"
	"github.com/spf13/viper"
)

func TestNew_WithFile_Success(t *testing.T) {
	tests := []struct {
		name          string
		givenOpts     []config.Option
		givenVar      string
		expectedValue string
	}{
		{
			name: "given a config file, expect values to be loaded from config",
			givenOpts: []config.Option{
				config.WithConfigFile("test.env"),
				config.WithEnvVars(false),
			},
			givenVar:      "TEST",
			expectedValue: "foobar",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := config.New(test.givenOpts...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			actual := viper.GetString(test.givenVar)

			if !cmp.Equal(actual, test.expectedValue) {
				t.Fatalf(cmp.Diff(actual, test.expectedValue))
			}
		})
	}
}

func TestNew_WithFile_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenOpts     []config.Option
		expectedValue error
	}{
		{
			name: "given an invalid config file, expect error to be returned",
			givenOpts: []config.Option{
				config.WithConfigFile("invalid.env"),
				config.WithEnvVars(false),
			},
			expectedValue: config.ErrFailedToLoadConfigFile,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := config.New(test.givenOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(err, test.expectedValue, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedValue, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestNew_WithEnvVar_Success(t *testing.T) {
	tests := []struct {
		name          string
		givenOpts     []config.Option
		givenEnvVars  map[string]string
		givenVar      string
		expectedValue string
	}{
		{
			name: "given a config file and env var, expect values from env var to be used",
			givenOpts: []config.Option{
				config.WithConfigFile("test.env"),
				config.WithEnvVars(true),
			},
			givenEnvVars: map[string]string{
				"TEST": "foobar_from_envvar",
			},
			givenVar:      "TEST",
			expectedValue: "foobar_from_envvar",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.givenEnvVars {
				t.Setenv(k, v)
			}

			s := os.Getenv("test")

			t.Log(s)

			_, err := config.New(test.givenOpts...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			actual := viper.GetString(test.givenVar)

			if !cmp.Equal(actual, test.expectedValue) {
				t.Fatalf(cmp.Diff(actual, test.expectedValue))
			}
		})
	}
}

func TestNew_WithEnvVar_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenOpts     []config.Option
		expectedValue error
	}{
		{
			name: "given an invalid config file, expect error to be returned",
			givenOpts: []config.Option{
				config.WithConfigFile("invalid.env"),
				config.WithEnvVars(false),
			},
			expectedValue: config.ErrFailedToLoadConfigFile,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := config.New(test.givenOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(err, test.expectedValue, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedValue, cmpopts.EquateErrors()))
			}
		})
	}
}
