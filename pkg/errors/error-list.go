package errors

import (
	"fmt"
)

var (
	ErrWarningNotification = ApiError{
		Code:   "WAR001",
		Title:  "The warning notification occured.",
		Status: "200",
	}

	// STATUS 400 - CODE: 'BRQXXX'
	ErrHeadersNotSupported = ApiError{
		Code: "BRQ001",
		Title: `The conditional headers provided in the request are not supported, 
		by the server.`,
		Status: "400",
	}

	ErrInvalidAuthenticationInfo = ApiError{
		Code: "BRQ002",
		Title: `The authentication information was not provided in the correct format. 
			Verify the value of Authorization header.`,
		Status: "400",
	}

	ErrInvalidHeaderValue = ApiError{
		Code:   "BRQ003",
		Title:  "The value provided in one of the HTTP headers was not in the correct format.",
		Status: "400",
	}

	ErrInvalidInput = ApiError{
		Code:   "BRQ004",
		Title:  "One of the request inputs is not valid.",
		Status: "400",
	}

	ErrInvalidQueryParameter = ApiError{
		Code:   "BRQ005",
		Title:  "An invalid value was specified for one of the query parameters in the request URI.",
		Status: "400",
	}

	ErrInvalidResourceName = ApiError{
		Code:   "BRQ006",
		Title:  "The specified resource name is not valid.",
		Status: "400",
	}

	ErrInvalidURI = ApiError{
		Code:   "BRQ007",
		Title:  "The requested URI does not represent any resource on the server.",
		Status: "400",
	}

	ErrInvalidJSONDocument = ApiError{
		Code:   "BRQ008",
		Title:  "The specified JSON is not syntatically valid.",
		Status: "400",
	}

	ErrInvalidJSONFieldValue = ApiError{
		Code:   "BRQ009",
		Title:  "The value provided for one of the JSON fields in the requested body was not in the correct format.",
		Status: "400",
	}

	ErrMD5Mismatch = ApiError{
		Code:   "BRQ010",
		Title:  "The MD5 value specified in the request did not match the MD5 value calculated by the server.",
		Status: "400",
	}

	ErrMetadataTooLarge = ApiError{
		Code:   "BRQ011",
		Title:  "The size of the specified metada exceeds the maximum size permitted.",
		Status: "400",
	}

	ErrMissingRequiredQueryParam = ApiError{
		Code:   "BRQ012",
		Title:  "A required query parameter was not specified for this request.",
		Status: "400",
	}

	ErrMissingRequiredHeader = ApiError{
		Code:   "BRQ013",
		Title:  "A required HTTP header was not specified.",
		Status: "400",
	}

	ErrMissingRequiredJSONField = ApiError{
		Code:   "BRQ014",
		Title:  "A required JSON field was not specified in the request body.",
		Status: "400",
	}

	ErrOutOfRangeInput = ApiError{
		Code:   "BRQ015",
		Title:  "One of the request inputs is out of range.",
		Status: "400",
	}

	ErrOutOfRangeQueryParameterValue = ApiError{
		Code:   "BRQ016",
		Title:  "A query parameter specified in the request URI is outside the permissible range.",
		Status: "400",
	}

	ErrUnsupportedHeader = ApiError{
		Code:   "BRQ017",
		Title:  "One of the HTTP headers specified in the request is not supported.",
		Status: "400",
	}

	ErrUnsupportedJSONField = ApiError{
		Code:   "BRQ018",
		Title:  "One of the JSON fields specified in the request body is not supported.",
		Status: "400",
	}

	ErrUnsupportedQueryParameter = ApiError{
		Code:   "BRQ019",
		Title:  "One of the query parameters in the request URI is not supported.",
		Status: "400",
	}

	ErrUnsupportedFilterOperator = ApiError{
		Code:   "BRQ020",
		Title:  "One of the filter operators is not supported.",
		Status: "400",
	}

	// STATUS 403, CODE: 'AUTHXX'
	ErrAccountDisabled = ApiError{
		Code:   "AUTH01",
		Title:  "The specified account is disabled.",
		Status: "403",
	}

	ErrAuthenticationFailed = ApiError{
		Code: "AUTH02",
		Title: `Server failed to authenticate the request. Make sure the value of 
		Authorization header is formed correctly including the signature.`,
		Status: "403",
	}

	ErrInsufficientAccPerm = ApiError{
		Code:   "AUTH03",
		Title:  "The account being accessed does not have sufficient permissions to execute this operation.",
		Status: "403",
	}
	ErrAuthInvalidCredentials = ApiError{
		Code:   "AUTH04",
		Title:  "Access is denied due to invalid credentials.",
		Status: "403",
	}

	ErrEndpointForbidden = ApiError{
		Code:   "FORB01",
		Title:  "Provided endpoint is forbidden.",
		Status: "403",
	}

	// STATUS 404, CODE: 'NTFXXX'
	ErrResourceNotFound = ApiError{
		Code:   "NTF001",
		Title:  "The specified resource does not exists.",
		Status: "404",
	}

	// STATUS 405, CODE: "MNAXXX"
	ErrMethodNotAllowed = ApiError{
		Code:   "MNA001",
		Title:  "The resource doesn't support the specified HTTP verb.",
		Status: "405",
	}

	// STATUS 406, CODE: "NALXXX"
	ErrLanguageNotAcceptable = ApiError{
		Code:   "NAL001",
		Title:  "The language provided within the json document is not supported.",
		Status: "406",
	}

	ErrLanguageHeaderNotAcceptable = ApiError{
		Code:   "NAL002",
		Title:  "The language provided in the request header is not supported.",
		Status: "406",
	}

	// STATUS 409, CODE: "CON001"
	ErrAccountAlreadyExists = ApiError{
		Code:   "CON001",
		Title:  "The Specified account already exists.",
		Status: "409",
	}

	ErrResourceAlreadyExists = ApiError{
		Code:   "CON002",
		Title:  "The specified resource already exists.",
		Status: "409",
	}

	// STATUS 413, CODE: 'RTLXXX'
	ErrRequestBodyTooLarge = ApiError{
		Code:   "RTL001",
		Title:  "The size of the request body exceeds the maximum size permitted.",
		Status: "413",
	}

	// STATUS 500, CODE: 'INTXXX'
	ErrInternalError = ApiError{
		Code:   "INT001",
		Title:  "The server encountered an internal error. Please retry the request.",
		Status: "500",
	}

	ErrOperatinTimedOut = ApiError{
		Code:   "INT002",
		Title:  "The operation could not be completed within the permitted time.",
		Status: "500",
	}

	// STATUS 503, CODE: 'UNAVXX'
	ErrServerBusy1 = ApiError{
		Code:   "UNAV01",
		Title:  "The server is currently unable to receive requests. Please retry your request.",
		Status: "503",
	}
	ErrServerBusy2 = ApiError{
		Code:   "UNAV02",
		Title:  "Operations per second is over the account limit.",
		Status: "503",
	}
)

func ErrTooManyNestedRelationships(relationship string) *ApiError {
	errObj := ErrUnsupportedQueryParameter.Copy()
	errObj.Detail = fmt.Sprintf(`Provided relationship: '%v', has to many nested relationships. 
		Only one level of nested relationships is allowed.`, relationship)
	return errObj
}