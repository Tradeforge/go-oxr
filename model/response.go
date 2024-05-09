package model

import (
	"encoding/json"
	"fmt"
)

// BaseResponse has all possible attributes that any response can use. It's intended to be embedded in a domain specific
// response struct.
type BaseResponse struct {
	PaginationHooks

	// Code is a code for the response status.
	Code int `json:"code,omitempty"`

	// A request id assigned by the server.
	RequestID string `json:"request_id,omitempty"`

	// A response message for successful requests.
	Message string `json:"message,omitempty"`

	// A map of error data for unsuccessful requests.
	ErrorData map[string]interface{} `json:"-"`
}

func (b *BaseResponse) UnmarshalJSON(data []byte) error {
	// We need to define an alias type to avoid infinite recursion when unmarshalling.
	type Alias BaseResponse
	bAlias := (*Alias)(b)
	if err := json.Unmarshal(data, bAlias); err != nil {
		return err
	}
	*b = BaseResponse(*bAlias)

	additionalParameters := make(map[string]interface{})
	if err := json.Unmarshal(data, &additionalParameters); err != nil {
		return err
	}
	delete(additionalParameters, "code")
	delete(additionalParameters, "request_id")
	delete(additionalParameters, "message")

	b.ErrorData = additionalParameters
	return nil
}

// PaginationHooks are links to next and/or previous pages. Embed this struct into an API response if the endpoint
// supports pagination.
type PaginationHooks struct {
	// If present, this value can be used to fetch the next page of data.
	NextURL string `json:"next_url,omitempty"`
}

func (p PaginationHooks) NextPage() string {
	return p.NextURL
}

// ResponseError represents an API response with an error status code.
type ResponseError struct {
	BaseResponse

	// An HTTP status code for unsuccessful requests.
	StatusCode int
}

// Error returns the details of an error response.
func (e *ResponseError) Error() string {
	return fmt.Sprintf("bad status with code '%d': message '%s': request ID '%s'", e.StatusCode, e.Message, e.RequestID)
}
