package mysql_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jamieaitken/cgs/mysql"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                string
		givenOpts           []mysql.Option
		expectedMaxLifetime time.Duration
		expectedMaxOpenCons int
		expectedMaxIdleCons int
	}{
		{
			name:                "given custom max lifetime, expect custom value and defaults for max open cons and idle cons",
			givenOpts:           []mysql.Option{mysql.WithMaxLifetime(time.Second * 120)},
			expectedMaxLifetime: time.Second * 120,
			expectedMaxIdleCons: 5,
			expectedMaxOpenCons: 0,
		},
		{
			name:                "given custom max open cons, expect custom value and defaults for max lifetime and idle cons",
			givenOpts:           []mysql.Option{mysql.WithMaxOpenConnections(20)},
			expectedMaxLifetime: time.Second * 60,
			expectedMaxIdleCons: 5,
			expectedMaxOpenCons: 20,
		},
		{
			name:                "given custom max idle cons, expect custom value and defaults for max lifetime and max open cons",
			givenOpts:           []mysql.Option{mysql.WithMaxIdleConnections(20)},
			expectedMaxLifetime: time.Second * 60,
			expectedMaxIdleCons: 20,
			expectedMaxOpenCons: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := mysql.New("root:hunter@(localhost:3306)/mysql?parseTime=true", test.givenOpts...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if !cmp.Equal(client.MaxLifetime(), test.expectedMaxLifetime) {
				t.Fatalf(cmp.Diff(client.MaxLifetime(), test.expectedMaxLifetime))
			}

			if !cmp.Equal(client.MaxIdleCons(), test.expectedMaxIdleCons) {
				t.Fatalf(cmp.Diff(client.MaxIdleCons(), test.expectedMaxIdleCons))
			}

			if !cmp.Equal(client.MaxOpenCons(), test.expectedMaxOpenCons) {
				t.Fatalf(cmp.Diff(client.MaxOpenCons(), test.expectedMaxOpenCons))
			}
		})
	}
}
