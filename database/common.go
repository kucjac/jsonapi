package database

import (
	"context"

	"github.com/neuronlabs/neuron/errors"
	"github.com/neuronlabs/neuron/log"
	"github.com/neuronlabs/neuron/mapping"
	"github.com/neuronlabs/neuron/query"
	"github.com/neuronlabs/neuron/repository"
)

// Count gets given scope models count.
func Count(ctx context.Context, db DB, s *query.Scope) (int64, error) {
	filterSoftDeleted(s)
	return getRepository(db, s).Count(ctx, s)
}

// Exists checks if given query model exists.
func Exists(ctx context.Context, db DB, s *query.Scope) (bool, error) {
	exister, isExister := getRepository(db, s).(repository.Exister)
	if !isExister {
		return false, errors.Wrapf(repository.ErrNotImplements, "repository for model: '%s' doesn't implement Exister interface", s.ModelStruct)
	}
	filterSoftDeleted(s)
	return exister.Exists(ctx, s)
}

func getRepository(db DB, s *query.Scope) repository.Repository {
	mapper := db.(repositoryMapper).mapper()
	repo, err := mapper.GetRepositoryByModelStruct(s.ModelStruct)
	if err != nil {
		log.Panic(err)
	}
	return repo
}

func getModelRepository(db DB, model *mapping.ModelStruct) repository.Repository {
	mapper := db.(repositoryMapper).mapper()
	repo, err := mapper.GetRepositoryByModelStruct(model)
	if err != nil {
		log.Panic(err)
	}
	return repo
}

func requireNoFilters(s *query.Scope) error {
	if len(s.Filters) != 0 {
		return errors.Wrapf(query.ErrInvalidInput, "given query doesn't allow filtering")
	}
	return nil
}

func errModelNotImplements(model *mapping.ModelStruct, interfaceName string) error {
	return errors.Wrapf(mapping.ErrModelNotImplements, "model: '%s' doesn't implement %s interface", model, interfaceName)
}

func logFormat(s *query.Scope, format string) string {
	return "SCOPE[" + s.ID.String() + "]" + format
}
