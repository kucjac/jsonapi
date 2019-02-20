package handler

import (
	"context"
	"github.com/kucjac/jsonapi/pkg/internal"
	ictrl "github.com/kucjac/jsonapi/pkg/internal/controller"
	"github.com/kucjac/jsonapi/pkg/internal/models"
	"github.com/kucjac/jsonapi/pkg/internal/query/paginations"
	iscope "github.com/kucjac/jsonapi/pkg/internal/query/scope"
	"github.com/kucjac/jsonapi/pkg/internal/query/sorts"
	"github.com/kucjac/jsonapi/pkg/log"
	"github.com/kucjac/jsonapi/pkg/mapping"
	"github.com/kucjac/jsonapi/pkg/query/scope"
	"github.com/kucjac/uni-db"
	"net/http"
)

// HandleList returns the handler function that handles 'List' operation for the 'm' model.
func (h *Handler) HandleList(m *mapping.ModelStruct) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Debugf("[LIST][%s] Begins", m.Collection())
		defer func() { log.Debugf("[LIST][%s] Finished", m.Collection()) }()
		// Add Flags to the scope setter

		ctx := context.WithValue(
			req.Context(),
			internal.ControllerIDCtxKey,
			h.c,
		)

		s, errs, err := (*ictrl.Controller)(h.c).QueryBuilder().BuildScopeMany(
			ctx,
			(*models.ModelStruct)(m),
			req.URL,
			h.c.Flags,
		)
		// handle internal error
		if err != nil {
			log.Errorf("Building Scope List failed: %v", err)
			h.internalError(rw)
			return
		}

		log.Debugf("Controller ptr: %p", s.Context().Value(internal.ControllerIDCtxKey))

		// handle client side errors
		if len(errs) > 0 {
			log.Debugf("BuildScopeList - ClientSide Errors: %v", errs)
			h.marshalErrors(rw, unsetStatus, errs...)
			return
		}

		/**

		TO DO:
		- set language filters

		*/

		if cfg := (*models.ModelStruct)(m).Config(); cfg != nil {
			eCfg := cfg.Endpoints.List
			if s.Pagination() == nil {
				// PresetPagination
				if eCfg.PresetPagination != nil && !eCfg.PresetPagination.IsZero() {
					p := paginations.NewFromConfig(eCfg.PresetPagination)

					err := iscope.SetPagination(s, p)
					if err != nil {
						log.Errorf("Preset Pagination for the model: '%v' has invalid config.", m.Type().String())
						h.internalError(rw)
						return
					}

				}
			}

			// PresetSorts
			if len(eCfg.PresetSorts) != 0 {
				var sortFields []*sorts.SortField
				for _, sort := range eCfg.PresetSorts {
					sortField, err := sorts.NewRawSortField((*models.ModelStruct)(m), sort)
					if err != nil {
						log.Errorf("Preset sort creation failed. Err: %v", err)
						h.internalError(rw)
						return
					}
					sortFields = append(sortFields, sortField)
				}

				// Append presetSortFields to the front of the sorts
				s.AppendSortFields(true, sortFields...)
			}
		}

		if s.Pagination() == nil && h.ListPagination != nil {
			if err := iscope.SetPagination(s, (*paginations.Pagination)(h.ListPagination)); err != nil {
				log.Errorf("Handler contains invalid default list pagination. %v", err)
				h.internalError(rw)
				return
			}
		}

		// List the values for the scope
		if err := (*scope.Scope)(s).List(); err != nil {
			var isNoResult bool
			if dbErr, ok := err.(*unidb.Error); ok {
				isNoResult = dbErr.Compare(unidb.ErrNoResult)
			}
			if !isNoResult {
				h.handleDBError(err, rw)
				return
			}
		}

		h.marshalScope(s, rw)
		return

	})
}
