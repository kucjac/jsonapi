package query

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/neuronlabs/neuron/config"
	"github.com/neuronlabs/neuron/controller"
	"github.com/neuronlabs/neuron/mapping"
)

// TestNewStringFilter tests the NewURLStringFilter function.
func TestNewUrlStringFilter(t *testing.T) {
	c := controller.NewDefault()

	err := c.RegisterRepository(repoName, &config.Repository{DriverName: repoName})
	require.NoError(t, err)

	require.NoError(t, c.RegisterModels(&TestingModel{}, &FilterRelationModel{}))

	mStruct, err := c.ModelStruct(&TestingModel{})
	require.NoError(t, err)

	t.Run("Primary", func(t *testing.T) {
		t.Run("WithoutOperator", func(t *testing.T) {
			filter, err := NewURLStringFilter(c, "filter[testing_models][id]", 521)
			require.NoError(t, err)

			assert.Equal(t, mStruct.Primary(), filter.StructField)
			require.Len(t, filter.Values, 1)

			fv := filter.Values[0]
			assert.Equal(t, OpEqual, fv.Operator)
			require.Len(t, fv.Values, 1)
			assert.Equal(t, 521, fv.Values[0])
		})

		t.Run("WithoutFilterWord", func(t *testing.T) {
			filter, err := NewURLStringFilter(c, "[testing_models][id][$ne]", "some string value")
			require.NoError(t, err)

			assert.Equal(t, mStruct.Primary(), filter.StructField)
			require.Len(t, filter.Values, 1)

			fv := filter.Values[0]
			assert.Equal(t, OpNotEqual, fv.Operator)
			require.Len(t, fv.Values, 1)
			assert.Equal(t, "some string value", fv.Values[0])
		})
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Run("NeuronCollectionName", func(t *testing.T) {
			_, err := NewURLStringFilter(c, "filter[invalid-collection][field_name][$eq]", 1)
			require.Error(t, err)
		})

		t.Run("Operator", func(t *testing.T) {
			_, err := NewURLStringFilter(c, "filter[testing_models][id][$unknown]", 1)
			require.Error(t, err)
		})

		t.Run("FieldName", func(t *testing.T) {
			_, err := NewURLStringFilter(c, "filter[testing_models][field-unknown][$eq]", "", 1)
			require.Error(t, err)
		})
	})

	t.Run("Attribute", func(t *testing.T) {
		filter, err := NewURLStringFilter(c, "[testing_models][attr][$ne]", "some string value")
		require.NoError(t, err)

		attrField, ok := mStruct.FieldByName("Attr")
		require.True(t, ok)

		assert.Equal(t, attrField, filter.StructField)
		require.Len(t, filter.Values, 1)

		fv := filter.Values[0]
		assert.Equal(t, OpNotEqual, fv.Operator)
		require.Len(t, fv.Values, 1)
		assert.Equal(t, "some string value", fv.Values[0])
	})

	t.Run("ForeignKey", func(t *testing.T) {
		_, err := NewURLStringFilter(c, "[testing_models][foreign_key][$ne]", "some string value")
		require.Error(t, err)

		filter, err := NewStringFilterWithForeignKey(c, "[testing_models][foreign_key][$ne]", "some string value")
		require.NoError(t, err)

		attrField, ok := mStruct.FieldByName("ForeignKey")
		require.True(t, ok)

		assert.Equal(t, attrField, filter.StructField)
		require.Len(t, filter.Values, 1)

		fv := filter.Values[0]
		assert.Equal(t, OpNotEqual, fv.Operator)
		require.Len(t, fv.Values, 1)
		assert.Equal(t, "some string value", fv.Values[0])
	})

	t.Run("Relationship", func(t *testing.T) {
		filter, err := NewURLStringFilter(c, "[testing_models][relation][id][$ne]", "some string value")
		require.NoError(t, err)

		relationField, ok := mStruct.RelationByName("Relation")
		require.True(t, ok)

		assert.Equal(t, relationField, filter.StructField)
		require.Len(t, filter.Nested, 1)

		nested := filter.Nested[0]
		require.Len(t, nested.Values, 1)

		fv := nested.Values[0]
		assert.Equal(t, OpNotEqual, fv.Operator)
		require.Len(t, fv.Values, 1)
		assert.Equal(t, "some string value", fv.Values[0])
	})
}

// TestFilterFormatQuery checks the FormatQuery function for the filters.
//noinspection GoNilness
func TestFilterFormatQuery(t *testing.T) {
	c := controller.NewDefault()

	err := c.RegisterRepository(repoName, &config.Repository{DriverName: repoName})
	require.NoError(t, err)

	require.NoError(t, c.RegisterModels(&TestingModel{}, &FilterRelationModel{}))

	mStruct, err := c.ModelStruct(&TestingModel{})
	require.NoError(t, err)

	t.Run("MultipleValue", func(t *testing.T) {
		tm := time.Now()
		f := NewFilterField(mStruct.Primary(), OpIn, 1, 2.01, 30, "something", []string{"i", "am"}, true, tm, &tm)
		q := f.FormatQuery()
		require.NotNil(t, q)

		assert.Len(t, q, 1)
		var k string
		var v []string

		for k, v = range q {
		}

		assert.Equal(t, fmt.Sprintf("filter[%s][%s][%s]", mStruct.Collection(), mStruct.Primary().NeuronName(), OpIn.URLAlias), k)

		if assert.Len(t, v, 1) {
			v = strings.Split(v[0], mapping.AnnotationSeparator)
			assert.Equal(t, "1", v[0])
			assert.Contains(t, v[1], "2.01")
			assert.Equal(t, "30", v[2])
			assert.Equal(t, "something", v[3])
			assert.Equal(t, "i", v[4])
			assert.Equal(t, "am", v[5])
			assert.Equal(t, "true", v[6])
			assert.Equal(t, fmt.Sprintf("%d", tm.Unix()), v[7])
			assert.Equal(t, fmt.Sprintf("%d", tm.Unix()), v[8])
		}
	})

	t.Run("WithNested", func(t *testing.T) {
		rel, ok := mStruct.RelationByName("relation")
		require.True(t, ok)

		relFilter := newRelationshipFilter(rel, NewFilterField(rel.ModelStruct().Primary(), OpIn, uint(1), uint64(2)))
		q := relFilter.FormatQuery()

		require.Len(t, q, 1)
		var k string
		var v []string

		for k, v = range q {
		}

		assert.Equal(t, fmt.Sprintf("filter[%s][%s][%s][%s]", mStruct.Collection(), relFilter.StructField.NeuronName(), relFilter.StructField.Relationship().Struct().Primary().NeuronName(), OpIn.URLAlias), k)
		if assert.Len(t, v, 1) {
			assert.NotNil(t, v)
			v = strings.Split(v[0], mapping.AnnotationSeparator)

			assert.Equal(t, "1", v[0])
			assert.Equal(t, "2", v[1])
		}
	})
}
