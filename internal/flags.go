package internal

import (
	"github.com/kucjac/jsonapi/flags"
)

// scopeCtxFlags - flags settable on scope, endpoint, model handler and controller
var ScopeCtxFlags []uint = []uint{
	flags.AddMetaCountList,
	flags.UseLinks,
	flags.ReturnPatchContent,
}

// modelCtxFlags - flags settable for model handler and controller
var ModelCtxFlags []uint = []uint{
	flags.AllowClientID,
}

// controllerContextFlags - flags settable only for the needs of the controller
var ControllerCtxFlags []uint = []uint{
	flags.AllowForeignKeyFilter,
	flags.UseFilterValueLimit,
}