package controller

import (
	"bytes"
	"github.com/kucjac/uni-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	buf := bytes.Buffer{}

	prepare := func(t *testing.T, models ...interface{}) {
		t.Helper()
		clearMap()
		if *debugFlag == true {
			basic := c.logger.(*unilogger.BasicLogger)
			basic.SetLevel(unilogger.DEBUG)
			c.logger = basic
		}
		buf.Reset()
		require.NoError(t, c.PrecomputeModels(models...))
	}

	prepareBlogs := func(t *testing.T) {
		prepare(t, &Blog{}, &Post{}, &Comment{})
	}

	tests := map[string]func(*testing.T){
		"single": func(t *testing.T) {
			prepareBlogs(t)

			value := &Blog{ID: 5, Title: "My title", ViewCount: 14}
			if assert.NoError(t, c.Marshal(&buf, value)) {
				marshaled := buf.String()
				assert.Contains(t, marshaled, `"title":"My title"`)
				assert.Contains(t, marshaled, `"view_count":14`)
				assert.Contains(t, marshaled, `"id":"5"`)
			}
		},
		"Time": func(t *testing.T) {
			type ModelPtrTime struct {
				ID   int        `jsonapi:"type=primary"`
				Time *time.Time `jsonapi:"type=attr"`
			}

			type ModelTime struct {
				ID   int       `jsonapi:"type=primary"`
				Time time.Time `jsonapi:"type=attr"`
			}

			t.Run("NoPtr", func(t *testing.T) {
				prepare(t, &ModelTime{})
				now := time.Now()
				v := &ModelTime{ID: 5, Time: now}
				if assert.NoError(t, c.Marshal(&buf, v)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, "time")
					assert.Contains(t, marshaled, `"id":"5"`)
				}
			})

			t.Run("Ptr", func(t *testing.T) {
				prepare(t, &ModelPtrTime{})
				now := time.Now()
				v := &ModelPtrTime{ID: 5, Time: &now}
				if assert.NoError(t, c.Marshal(&buf, v)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, "time")
					assert.Contains(t, marshaled, `"id":"5"`)
				}
			})

		},
		"singleWithMap": func(t *testing.T) {
			t.Run("PtrString", func(t *testing.T) {
				type MpString struct {
					ID  int                `jsonapi:"type=primary"`
					Map map[string]*string `jsonapi:"type=attr"`
				}
				prepare(t, &MpString{})

				kv := "some"
				value := &MpString{ID: 5, Map: map[string]*string{"key": &kv}}
				if assert.NoError(t, c.Marshal(&buf, value)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"map":{"key":"some"}`)
				}
			})

			t.Run("NilString", func(t *testing.T) {
				type MpString struct {
					ID  int                `jsonapi:"type=primary"`
					Map map[string]*string `jsonapi:"type=attr"`
				}
				prepare(t, &MpString{})
				value := &MpString{ID: 5, Map: map[string]*string{"key": nil}}
				if assert.NoError(t, c.Marshal(&buf, value)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"map":{"key":null}`)
				}
			})

			t.Run("PtrInt", func(t *testing.T) {
				type MpInt struct {
					ID  int             `jsonapi:"type=primary"`
					Map map[string]*int `jsonapi:"type=attr"`
				}
				prepare(t, &MpInt{})

				kv := 5
				value := &MpInt{ID: 5, Map: map[string]*int{"key": &kv}}
				if assert.NoError(t, c.Marshal(&buf, value)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"map":{"key":5}`)
				}
			})
			t.Run("NilPtrInt", func(t *testing.T) {
				type MpInt struct {
					ID  int             `jsonapi:"type=primary"`
					Map map[string]*int `jsonapi:"type=attr"`
				}
				prepare(t, &MpInt{})

				value := &MpInt{ID: 5, Map: map[string]*int{"key": nil}}
				if assert.NoError(t, c.Marshal(&buf, value)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"map":{"key":null}`)
				}
			})
			t.Run("PtrFloat", func(t *testing.T) {
				type MpFloat struct {
					ID  int                 `jsonapi:"type=primary"`
					Map map[string]*float64 `jsonapi:"type=attr"`
				}
				prepare(t, &MpFloat{})

				fv := 1.214
				value := &MpFloat{ID: 5, Map: map[string]*float64{"key": &fv}}
				if assert.NoError(t, c.Marshal(&buf, value)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"map":{"key":1.214}`)
				}
			})
			t.Run("NilPtrFloat", func(t *testing.T) {
				type MpFloat struct {
					ID  int                 `jsonapi:"type=primary"`
					Map map[string]*float64 `jsonapi:"type=attr"`
				}
				prepare(t, &MpFloat{})

				value := &MpFloat{ID: 5, Map: map[string]*float64{"key": nil}}
				if assert.NoError(t, c.Marshal(&buf, value)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"map":{"key":null}`)
				}
			})

			t.Run("SliceInt", func(t *testing.T) {
				type MpSliceInt struct {
					ID  int              `jsonapi:"type=primary"`
					Map map[string][]int `jsonapi:"type=attr"`
				}
				prepare(t, &MpSliceInt{})

				value := &MpSliceInt{ID: 5, Map: map[string][]int{"key": {1, 5}}}
				if assert.NoError(t, c.Marshal(&buf, value)) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"map":{"key":[1,5]}`)
				}
			})

		},
		"many": func(t *testing.T) {
			prepareBlogs(t)

			values := []*Blog{{ID: 5, Title: "First"}, {ID: 2, Title: "Second"}}
			if assert.NoError(t, c.Marshal(&buf, values)) {
				marshaled := buf.String()
				assert.Contains(t, marshaled, `"title":"First"`)
				assert.Contains(t, marshaled, `"title":"Second"`)

				assert.Contains(t, marshaled, `"id":"5"`)
				assert.Contains(t, marshaled, `"id":"2"`)
				t.Log(marshaled)
			}
		},
		"Nested": func(t *testing.T) {

			t.Run("Simple", func(t *testing.T) {
				type NestedSub struct {
					First int
				}

				type Simple struct {
					ID     int        `jsonapi:"type=primary"`
					Nested *NestedSub `jsonapi:"type=attr"`
				}

				prepare(t, &Simple{})

				err := c.Marshal(&buf, &Simple{ID: 2, Nested: &NestedSub{First: 1}})
				if assert.NoError(t, err) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"nested":{"first":1}`)
				}
			})

			t.Run("DoubleNested", func(t *testing.T) {

				type NestedSub struct {
					First int
				}

				type DoubleNested struct {
					Nested *NestedSub
				}

				type Simple struct {
					ID     int           `jsonapi:"type=primary"`
					Double *DoubleNested `jsonapi:"type=attr"`
				}

				prepare(t, &Simple{})

				err := c.Marshal(&buf, &Simple{ID: 2, Double: &DoubleNested{Nested: &NestedSub{First: 1}}})
				if assert.NoError(t, err) {
					marshaled := buf.String()
					assert.Contains(t, marshaled, `"nested":{"first":1}`)
					assert.Contains(t, marshaled, `"double":{"nested"`)
				}
			})

		},
	}

	for name, testFunc := range tests {

		t.Run(name, testFunc)
	}

}

func TestMarshalScope(t *testing.T) {

	buf := bytes.NewBufferString("")

	c.PrecomputeModels(&Blog{}, &Post{}, &Comment{})

	req := httptest.NewRequest("GET", `/blogs/3?include=posts,current_post.latest_comment&fields[blogs]=title,created_at,posts&fields[posts]=title,body,comments`, nil)
	scope, errs, err := c.BuildScopeSingle(req, &Endpoint{Type: Get}, &ModelHandler{ModelType: reflect.TypeOf(Blog{})})
	assertNil(t, err)
	assertEmpty(t, errs)
	scope.Value = &Blog{ID: 3, Title: "My own title.", CreatedAt: time.Now(), Posts: []*Post{{ID: 1}}, CurrentPost: &Post{ID: 2}}

	// assertEqual(t, 1, len(scope.IncludedFields))
	postInclude := scope.IncludedFields[0]
	postScope := postInclude.Scope
	postScope.Value = []*Post{{ID: 1, Title: "Post title", Body: "Post body."}}

	currentPost := scope.IncludedFields[1]
	currentPost.Scope.Value = &Post{ID: 2, Title: "Current One", Body: "This is current post", LatestComment: &Comment{ID: 1}}

	latestComment := currentPost.Scope.IncludedFields[0]
	latestComment.Scope.Value = &Comment{ID: 1, Body: "This is such a great post", PostID: 2}

	err = scope.SetCollectionValues()
	assertNil(t, err)
	for scope.NextIncludedField() {
		includedField, err := scope.CurrentIncludedField()
		assertNil(t, err)
		_, err = includedField.GetMissingPrimaries()
		assertNil(t, err)
		err = includedField.Scope.SetCollectionValues()
		assertNil(t, err)

		for includedField.Scope.NextIncludedField() {
			includedField.Scope.CurrentIncludedField()
			nestedIncluded, err := includedField.Scope.CurrentIncludedField()
			assertNil(t, err)
			_, err = nestedIncluded.GetMissingPrimaries()
			assertNil(t, err)

			err = nestedIncluded.Scope.SetCollectionValues()
			assertNil(t, err)
		}
	}

	payload, err := marshalScope(scope, c)
	assertNoError(t, err)

	err = MarshalPayload(buf, payload)
	assertNoError(t, err)
	// even if included, there is no

	assertTrue(t, strings.Contains(buf.String(), "{\"type\":\"posts\",\"id\":\"1"))
	assertTrue(t, strings.Contains(buf.String(), "\"title\":\"My own title.\""))
	assertTrue(t, strings.Contains(buf.String(), "{\"type\":\"comments\",\"id\":\"1\",\"attributes\":{\"body\""))
	// assertTrue(t, strings.Contains(buf.String(), "\"relationships\":{\"posts\":{\"data\":[{\"type\":\"posts\",\"id\":\"1\"}]}}"))
	assertTrue(t, strings.Contains(buf.String(), "\"type\":\"blogs\",\"id\":\"3\""))
	clearMap()
	buf.Reset()

	scope = getBlogScope()
	errs = scope.buildIncludeList("current_post")
	assertEmpty(t, errs)
	scope.Value = &Blog{ID: 4, Title: "The title.", CreatedAt: time.Now(), CurrentPost: &Post{ID: 3}}
	errs = scope.buildFieldset("title", "created_at", "current_post")
	assertEmpty(t, errs)

	scope.IncludedScopes[c.MustGetModelStruct(&Post{})].Value = &Post{ID: 3, Title: "Breaking News!", Body: "Some body"}
	errs = scope.IncludedScopes[c.MustGetModelStruct(&Post{})].buildFieldset("title", "body")
	assertEmpty(t, errs)

	payload, err = marshalScope(scope, c)
	assertNil(t, err)
	err = MarshalPayload(buf, payload)
	assertNil(t, err)

	// assertTrue(t, strings.Contains(buf.String(),
	// "\"relationships\":{\"current_post\":{\"data\":{\"type\":\"posts\",\"id\":\"3\"}}}"))

	// t.Log(buf.String())
	assertTrue(t, strings.Contains(buf.String(),
		"\"type\":\"blogs\",\"id\":\"4\",\"attributes\":{\"created_at\":"))
	// assertTrue(t, strings.Contains(buf.String(),
	// "\"included\":[{\"type\":\"posts\",\"id\":\"3\",\"attributes\":{\"body\":\"Some body\",\"title\":\"Breaking News!\"}}]"))

	clearMap()
	buf.Reset()
	scope = getBlogScope()
	scope.Value = []*Blog{{ID: 4, Title: "The title one."}, {ID: 5, Title: "The title two"}}
	errs = scope.buildFieldset("title")
	assertEmpty(t, errs)

	payload, err = marshalScope(scope, c)
	assertNil(t, err)

	err = MarshalPayload(buf, payload)
	assertNil(t, err)

	assertTrue(t, strings.Contains(buf.String(), `"type":"blogs","id":"4","attributes":{"title":"The title one."}`))
	assertTrue(t, strings.Contains(buf.String(), `"type":"blogs","id":"5","attributes":{"title":"The title two"}`))

	// scope with no value
	clearMap()
	buf.Reset()
	scope = getBlogScope()

	payload, err = marshalScope(scope, c)
	assertError(t, err)

	err = MarshalPayload(buf, payload)
	assertNil(t, err)

	t.Run("MarshalToManyRelationship", func(t *testing.T) {
		clearMap()
		require.NoError(t, c.PrecomputeModels(&Pet{}, &User{}))

		scope, err := c.NewScope(&Pet{})
		require.NoError(t, err)

		scope.Value = &Pet{ID: 5, Owners: []*User{{ID: 2}, {ID: 3}}}
		scope.SetFields("Owners")

		payload, err := c.MarshalScope(scope)
		if assert.NoError(t, err) {
			single, ok := payload.(*OnePayload)
			if assert.True(t, ok) {
				if assert.NotNil(t, single.Data) {
					if assert.NotEmpty(t, single.Data.Relationships) {
						if assert.NotNil(t, single.Data.Relationships["owners"]) {
							owners, ok := single.Data.Relationships["owners"].(*RelationshipManyNode)
							if assert.True(t, ok) {
								var count int
								for _, owner := range owners.Data {
									if assert.NotNil(t, owner) {
										switch owner.ID {
										case "2", "3":
											count += 1
										}
									}
								}
								assert.Equal(t, 2, count)
							}

						}
					}
				}
			}
		}

	})

	t.Run("MarshalToManyEmptyRelationship", func(t *testing.T) {
		clearMap()
		require.NoError(t, c.PrecomputeModels(&Pet{}, &User{}))

		scope, err := c.NewScope(&Pet{})
		require.NoError(t, err)

		scope.Value = &Pet{ID: 5, Owners: []*User{}}
		scope.SetFields("Owners")

		payload, err := c.MarshalScope(scope)
		if assert.NoError(t, err) {
			single, ok := payload.(*OnePayload)
			if assert.True(t, ok) {
				if assert.NotNil(t, single.Data) {
					if assert.NotEmpty(t, single.Data.Relationships) {
						if assert.NotNil(t, single.Data.Relationships["owners"]) {
							owners, ok := single.Data.Relationships["owners"].(*RelationshipManyNode)
							if assert.True(t, ok, reflect.TypeOf(single.Data.Relationships["owners"]).String()) {
								if assert.NotNil(t, owners) {
									assert.Empty(t, owners.Data)
								}
							}

						}
					}
				}
				buf := bytes.Buffer{}
				assert.NoError(t, MarshalPayload(&buf, single))
				assert.Contains(t, buf.String(), "owners")
			}
		}

	})
}

func TestMarshalScopeRelationship(t *testing.T) {
	clearMap()
	getBlogScope()
	req := httptest.NewRequest("GET", "/blogs/1/relationships/posts", nil)
	scope, errs, err := c.BuildScopeRelationship(req, &Endpoint{Type: GetRelationship}, &ModelHandler{ModelType: reflect.TypeOf(Blog{})})

	assertNil(t, err)
	assertEmpty(t, errs)

	scope.Value = &Blog{ID: 1, Posts: []*Post{{ID: 1}, {ID: 3}}}

	postsScope, err := scope.GetRelationshipScope()
	assertNil(t, err)

	payload, err := c.MarshalScope(postsScope)
	assertNil(t, err)

	buffer := bytes.NewBufferString("")

	err = MarshalPayload(buffer, payload)
	assertNil(t, err)

}

type HiddenModel struct {
	ID          int    `jsonapi:"type=primary;flags=hidden"`
	Visibile    string `jsonapi:"type=attr"`
	HiddenField string `jsonapi:"type=attr;flags=hidden"`
}

func (h *HiddenModel) CollectionName() string {
	return "hiddens"
}

func TestMarshalHiddenScope(t *testing.T) {

	clearMap()
	assertNoError(t, c.PrecomputeModels(&HiddenModel{}), failNow)

	scope, err := c.NewScope(&HiddenModel{})
	assertNoError(t, err, failNow)

	scope.Value = &HiddenModel{ID: 1, Visibile: "Visible", HiddenField: "Invisible"}

	payload, err := c.MarshalScope(scope)
	assertNoError(t, err, failNow)

	buffer := bytes.NewBufferString("")
	err = MarshalPayload(buffer, payload)
	assertNoError(t, err, failNow)

}