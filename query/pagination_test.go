package query

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatQuery tests the FormatQuery method.
func TestFormatQuery(t *testing.T) {
	t.Run("Paged", func(t *testing.T) {
		p := PagedPagination(10, 2)
		require.NoError(t, p.Check())

		q := url.Values{}

		p.FormatQuery(q)
		require.Len(t, q, 2)

		assert.Equal(t, "10", q.Get(ParamPageNumber))
		assert.Equal(t, "2", q.Get(ParamPageSize))
	})

	t.Run("Limited", func(t *testing.T) {
		p := LimitOffsetPagination(10, 0)

		require.NoError(t, p.Check())

		q := p.FormatQuery()
		require.Len(t, q, 1)

		assert.Equal(t, "10", q.Get(ParamPageLimit))
	})

	t.Run("Offseted", func(t *testing.T) {
		p := LimitOffsetPagination(0, 10)

		require.NoError(t, p.Check())

		q := p.FormatQuery()
		require.Len(t, q, 1)

		assert.Equal(t, "10", q.Get(ParamPageOffset))
	})

	t.Run("LimitOffset", func(t *testing.T) {
		p := LimitOffsetPagination(10, 140)

		require.NoError(t, p.Check())

		q := p.FormatQuery()
		require.Len(t, q, 2)

		assert.Equal(t, "10", q.Get(ParamPageLimit))
		assert.Equal(t, "140", q.Get(ParamPageOffset))
	})
}
