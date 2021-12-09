package redis_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jamieaitken/cgs/redis"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name               string
		givenAddrs         []string
		givenOpts          []redis.Option
		expectedAddrs      []string
		expectedMasterName string
		expectedPassword   string
		expectedDB         int
	}{
		{
			name:               "given custom master name, expect defaults to be used for rest",
			givenAddrs:         []string{"test"},
			givenOpts:          []redis.Option{redis.WithMasterName("test")},
			expectedAddrs:      []string{"test"},
			expectedMasterName: "test",
			expectedDB:         0,
			expectedPassword:   "",
		},
		{
			name:               "given custom password, expect defaults to be used for rest",
			givenAddrs:         []string{"test"},
			givenOpts:          []redis.Option{redis.WithPassword("test")},
			expectedAddrs:      []string{"test"},
			expectedMasterName: "mymaster",
			expectedDB:         0,
			expectedPassword:   "test",
		},
		{
			name:               "given custom db, expect defaults to be used for rest",
			givenAddrs:         []string{"test"},
			givenOpts:          []redis.Option{redis.WithDB(1)},
			expectedAddrs:      []string{"test"},
			expectedMasterName: "mymaster",
			expectedDB:         1,
			expectedPassword:   "",
		},
		{
			name:               "given custom client func, expect defaults to be used for rest",
			givenAddrs:         []string{"test"},
			givenOpts:          []redis.Option{redis.WithClientType(redis.NonFailOver)},
			expectedAddrs:      []string{"test"},
			expectedMasterName: "mymaster",
			expectedDB:         0,
			expectedPassword:   "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := redis.New(test.givenAddrs, test.givenOpts...)

			if !cmp.Equal(actual.Addrs(), test.expectedAddrs) {
				t.Fatalf(cmp.Diff(actual.Addrs(), test.expectedAddrs))
			}

			if !cmp.Equal(actual.MasterName(), test.expectedMasterName) {
				t.Fatalf(cmp.Diff(actual.MasterName(), test.expectedMasterName))
			}

			if !cmp.Equal(actual.Password(), test.expectedPassword) {
				t.Fatalf(cmp.Diff(actual.Password(), test.expectedPassword))
			}

			if !cmp.Equal(actual.DB(), test.expectedDB) {
				t.Fatalf(cmp.Diff(actual.DB(), test.expectedDB))
			}
		})
	}
}
