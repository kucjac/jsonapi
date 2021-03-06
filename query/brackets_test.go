package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSplitBracketParameter tests the SplitBracketParameter function.
func TestSplitBracketParameter(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCase := "[collection][field][$operator]"
		split, err := SplitBracketParameter(testCase)
		require.NoError(t, err)

		if assert.Len(t, split, 3) {
			assert.Equal(t, "collection", split[0])
			assert.Equal(t, "field", split[1])
			assert.Equal(t, "$operator", split[2])
		}
	})

	t.Run("DoubleOpen", func(t *testing.T) {
		testCase := "[[collection][field][$operator]"
		_, err := SplitBracketParameter(testCase)
		require.Error(t, err)
	})

	t.Run("DoubleClose", func(t *testing.T) {
		testCase := "[collection]][field][$operator]"
		_, err := SplitBracketParameter(testCase)
		require.Error(t, err)
	})

	t.Run("NoClose", func(t *testing.T) {
		testCase := "[collection"
		_, err := SplitBracketParameter(testCase)
		require.Error(t, err)
	})

	t.Run("Multiple", func(t *testing.T) {
		type stringBool struct {
			Str string
			Val bool
		}

		values := []stringBool{
			{"[some][thing]", true},
			{"[no][closing", false},
			{"no][opening]", false},
			{"]justclosing", false},
			{"[doubleopen[]", false},
			{"[doubleclose]]", false},
		}

		var split []string
		var err error
		for _, v := range values {
			split, err = SplitBracketParameter(v.Str)
			if !v.Val {
				assert.Error(t, err)
				// t.Log(err)
				if err == nil {
					t.Log(v.Str)
				}
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, split)
			}
		}
	})
}
