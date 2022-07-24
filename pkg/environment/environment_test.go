package environment

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStringDefaultValue(t *testing.T) {
	t.Parallel()
	c := require.New(t)

	defaultValue := "val"
	c.Equal(defaultValue, GetString("GET_STRING", defaultValue))
}

func TestGetStringCustomValue(t *testing.T) {
	t.Parallel()
	c := require.New(t)

	withTestEnv("custom", func(varName string) {
		c.Equal("custom", GetString(varName, ""))
	})
}

func TestGetInt64DefaultValue(t *testing.T) {
	t.Parallel()
	c := require.New(t)

	defaultValue := int64(872)
	c.Equal(defaultValue, GetInt64("GET_INT64", defaultValue))
}

func TestGetInt64CustomValue(t *testing.T) {
	t.Parallel()
	c := require.New(t)

	ttable := []struct {
		val        interface{}
		defaultVal int64
		expected   int64
	}{
		{false, -1, -1},
		{"34", -1, 34},
		{-34, -1, -34},
		{"00001", -1, 1},
		{"x", -1, -1},
		{0x7F, -1, 127},
	}

	for _, test := range ttable {
		withTestEnv(test.val, func(varName string) {
			c.Equal(test.expected, GetInt64(varName, test.defaultVal), "Failed to get %#v with %#v", varName, test.val)
		})
	}
}

func withTestEnv(val interface{}, codeBlock func(varName string)) {
	varName := fmt.Sprintf("TEST_%d_%d", rand.Intn(math.MaxInt32), rand.Intn(math.MaxInt32))
	_ = os.Setenv(varName, fmt.Sprintf("%v", val))

	defer func() {
		_ = os.Unsetenv(varName)
	}()

	codeBlock(varName)
}
