package salesforce

// Err implements the error interface so we can have constant errors.
type Err string

func (e Err) Error() string {
	return string(e)
}

// Error Constants
// Salesforce documents these as the error responses they will emit.
// For more detail see their docs: https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/errorcodes.htm
const (
	ErrMultipleExternalIDMatch  = Err("salesforce: external ID exists in more than one record")                      // 300
	ErrRequestContentNotChanged = Err("salesforce: request content has not changed since a specified date and time") // 304
	ErrBadRequest               = Err("salesforce: bad request")
	ErrUnauthorized             = Err("salesforce: unauthorized request")
	ErrForbidden                = Err("salesforce: forbidden")
	ErrMethodNotAllowed         = Err("salesforce: method not allowed")                              // 405
	ErrConflict                 = Err("salesforce: conflict with the current state of the resource") // 409
	ErrInternalError            = Err("salesforce: internal error")
	ErrUnknown                  = Err("salesforce: unexpected error occurred")
)

// BadRequestError represents the response sent by salesforce for a Bad Request 400 error
type BadRequestError struct {
	Message   string   `json:"message"`
	ErrorCode string   `json:"errorCode"`
	Fields    []string `json:"fields"`
}
