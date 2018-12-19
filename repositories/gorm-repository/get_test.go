package gormrepo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	tests := map[string]func(*testing.T){
		"AttrOnly": func(t *testing.T) {
			models := []interface{}{&Simple{}}
			c, err := prepareJSONAPI(models...)
			require.NoError(t, err)

			repo, err := prepareGORMRepo(models...)
			require.NoError(t, err)

			modelToTest := &Simple{ID: 1, Attr1: "First", Attr2: 2}
			require.NoError(t, repo.db.Create(modelToTest).Error)

			scope, err := c.NewScope(&Simple{})
			require.NoError(t, err)

			scope.Value = &Simple{}
			scope.SetPrimaryFilters(modelToTest.ID)

			scope.SetAllFields()

			err = repo.Get(scope)
			if assert.NoError(t, err) {
				simple, ok := scope.Value.(*Simple)
				require.True(t, ok)

				assert.Equal(t, modelToTest.ID, simple.ID)
				assert.Equal(t, modelToTest.Attr1, simple.Attr1)
				assert.Equal(t, modelToTest.Attr2, simple.Attr2)
			}
		},
		"RelationBelongsTo": func(t *testing.T) {
			models := []interface{}{&Comment{}, &Post{}}
			c, err := prepareJSONAPI(models...)
			require.NoError(t, err)

			repo, err := prepareGORMRepo(models...)
			require.NoError(t, err)

			post := &Post{ID: 1, Lang: "pl", Title: "Title"}
			require.NoError(t, repo.db.Create(post).Error)

			comment := &Comment{PostID: post.ID}
			require.NoError(t, repo.db.Create(comment).Error)

			scope, err := c.NewScope(&Comment{})
			require.NoError(t, err)

			scope.NewValueSingle()

			require.NoError(t, scope.SetFields(scope.Struct.GetPrimaryField(), "PostID"))

			scope.SetPrimaryFilters(comment.ID)

			err = repo.Get(scope)
			if assert.NoError(t, err) {
				comm, ok := scope.Value.(*Comment)
				require.True(t, ok)

				assert.Equal(t, comment.PostID, comm.PostID)
				assert.Equal(t, comment.ID, comm.ID)
			}
		},
		"RelationHasOne": func(t *testing.T) {
			t.Run("Synced", func(t *testing.T) {
				models := []interface{}{&Human{}, &BodyPart{}}
				c, err := prepareJSONAPI(models...)
				require.NoError(t, err)

				repo, err := prepareGORMRepo(models...)
				require.NoError(t, err)

				scope, err := c.NewScope(&Human{})
				require.NoError(t, err)

				human := &Human{ID: 4}

				nose := &BodyPart{ID: 5, HumanID: human.ID}
				require.NoError(t, repo.db.Create(human).Error)
				require.NoError(t, repo.db.Create(nose).Error)

				scope.Value = &Human{ID: human.ID}
				scope.SetPrimaryFilters(human.ID)

				require.NoError(t, scope.SetFields("Nose"))

				err = repo.Get(scope)
				if assert.NoError(t, err) {
					hum, ok := scope.Value.(*Human)
					require.True(t, ok)

					assert.Zero(t, hum.Nose)
				}

			})

			t.Run("NonSynced", func(t *testing.T) {
				models := []interface{}{&Human{}, &BodyPart{}}
				c, err := prepareJSONAPI(models...)
				require.NoError(t, err)

				repo, err := prepareGORMRepo(models...)
				require.NoError(t, err)

				scope, err := c.NewScope(&Human{})
				require.NoError(t, err)

				human := &Human{ID: 4}

				nose := &BodyPart{ID: 5, HumanNonSyncID: human.ID}
				require.NoError(t, repo.db.Create(human).Error)
				require.NoError(t, repo.db.Create(nose).Error)

				scope.Value = &Human{}

				scope.SetPrimaryFilters(human.ID)
				require.NoError(t, scope.SetFields("NoseNonSynced"))

				err = repo.Get(scope)
				if assert.NoError(t, err) {
					hum, ok := scope.Value.(*Human)
					require.True(t, ok)

					assert.Equal(t, human.ID, hum.ID)
					assert.Zero(t, hum.Nose)
					if assert.NotZero(t, hum.NoseNonSynced) {
						assert.Equal(t, nose.ID, hum.NoseNonSynced.ID)
					}
				}
			})
		},
		"RelationHasMany": func(t *testing.T) {
			t.Run("Synced", func(t *testing.T) {
				models := []interface{}{&Human{}, &BodyPart{}}
				c, err := prepareJSONAPI(models...)
				require.NoError(t, err)

				repo, err := prepareGORMRepo(models...)
				require.NoError(t, err)

				scope, err := c.NewScope(&Human{})
				require.NoError(t, err)

				human := &Human{ID: 4}

				earID1, earID2 := 5, 6
				ears := []*BodyPart{{ID: earID1, HumanID: human.ID}, {ID: earID2, HumanNonSyncID: human.ID}}
				require.NoError(t, repo.db.Create(human).Error)
				for _, ear := range ears {
					require.NoError(t, repo.db.Create(ear).Error)
				}

				scope.Value = &Human{}

				scope.SetPrimaryFilters(human.ID)
				require.NoError(t, scope.SetFields("Ears"))

				err = repo.Get(scope)
				if assert.NoError(t, err) {
					hum, ok := scope.Value.(*Human)
					require.True(t, ok)

					assert.Zero(t, hum.Ears)
				}

			})
			t.Run("NonSynced", func(t *testing.T) {
				models := []interface{}{&Human{}, &BodyPart{}}
				c, err := prepareJSONAPI(models...)
				require.NoError(t, err)

				repo, err := prepareGORMRepo(models...)
				require.NoError(t, err)

				scope, err := c.NewScope(&Human{})
				require.NoError(t, err)

				human := &Human{ID: 4}

				earID1, earID2 := 5, 6
				ears := []*BodyPart{{ID: earID1, HumanNonSyncID: human.ID}, {ID: earID2, HumanNonSyncID: human.ID}}
				require.NoError(t, repo.db.Create(human).Error)
				for _, ear := range ears {
					require.NoError(t, repo.db.Create(ear).Error)
				}

				scope.Value = &Human{}

				scope.SetPrimaryFilters(human.ID)
				require.NoError(t, scope.SetFields("EarsNonSync"))

				err = repo.Get(scope)
				if assert.NoError(t, err) {
					hum, ok := scope.Value.(*Human)
					require.True(t, ok)

					assert.Equal(t, human.ID, hum.ID)
					assert.Zero(t, hum.Ears)
					if assert.NotZero(t, hum.EarsNonSync) {
						var count int
						for _, earNS := range hum.EarsNonSync {
							switch earNS.ID {
							case earID1, earID2:
								count++
							default:
								t.FailNow()
							}
						}
						assert.Equal(t, 2, count)
					}
				}

			})
		},
		"RelationMany2Many": func(t *testing.T) {
			t.Run("Synced", func(t *testing.T) {
				models := []interface{}{&M2MFirst{}, &M2MSecond{}}
				c, err := prepareJSONAPI(models...)
				require.NoError(t, err)

				repo, err := prepareGORMRepo(models...)
				require.NoError(t, err)

				first := &M2MFirst{}
				require.NoError(t, repo.db.Create(first).Error)
				secondOne := &M2MSecond{}
				secondTwo := &M2MSecond{}

				require.NoError(t, repo.db.Create(secondOne).Error)
				require.NoError(t, repo.db.Create(secondTwo).Error)

				require.NoError(t, repo.db.Model(first).Updates(&M2MFirst{ID: first.ID, SecondsSync: []*M2MSecond{secondOne, secondTwo}}).Error)

				scope, err := c.NewScope(&M2MFirst{})
				require.NoError(t, err)
				scope.NewValueSingle()

				require.NoError(t, scope.SetFields("SecondsSync"))

				if assert.NoError(t, repo.Get(scope)) {
					st, ok := scope.Value.(*M2MFirst)
					require.True(t, ok)

					assert.Empty(t, st.Seconds)
					assert.Empty(t, st.SecondsSync)

				}
			})
			t.Run("NonSynced", func(t *testing.T) {
				models := []interface{}{&M2MFirst{}, &M2MSecond{}}
				c, err := prepareJSONAPI(models...)
				require.NoError(t, err)

				repo, err := prepareGORMRepo(models...)
				require.NoError(t, err)

				first := &M2MFirst{}
				require.NoError(t, repo.db.Create(first).Error)
				secondOne := &M2MSecond{}
				secondTwo := &M2MSecond{}

				require.NoError(t, repo.db.Create(secondOne).Error)
				require.NoError(t, repo.db.Create(secondTwo).Error)

				require.NoError(t, repo.db.Model(first).Updates(&M2MFirst{ID: first.ID, Seconds: []*M2MSecond{secondOne, secondTwo}}).Error)

				scope, err := c.NewScope(&M2MFirst{})
				require.NoError(t, err)
				scope.NewValueSingle()

				require.NoError(t, scope.SetFields("Seconds"))

				if assert.NoError(t, repo.Get(scope)) {
					st, ok := scope.Value.(*M2MFirst)
					require.True(t, ok)

					assert.Len(t, st.Seconds, 2)
					assert.Empty(t, st.SecondsSync)

				}
			})
		},
		"Error": func(t *testing.T) {

		},
	}

	for name, testFunc := range tests {
		t.Run(name, testFunc)
	}
}