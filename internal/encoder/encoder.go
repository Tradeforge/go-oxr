package encoder

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Encoder defines a path and query param encoder that plays nicely with the Polygon REST API.
type Encoder struct {
	validate     *validator.Validate
	pathEncoder  *form.Encoder
	queryEncoder *form.Encoder
}

// New returns a new path and query param encoder.
func New() *Encoder {
	return &Encoder{
		validate:     validator.New(),
		pathEncoder:  newEncoder("path"),
		queryEncoder: newEncoder("query"),
	}
}

// EncodeParams encodes path and query params and returns a valid request URI.
func (e *Encoder) EncodeParams(path string, params any) (string, error) {
	if err := e.validateParams(params); err != nil {
		return path, err
	}

	uri, err := e.encodePath(path, params)
	if err != nil {
		return "", err
	}

	query, err := e.encodeQuery(params)
	if err != nil {
		return "", err
	} else if query != "" {
		uri += "?" + query
	}

	return uri, nil
}

func (e *Encoder) validateParams(params any) error {
	if params == nil {
		return nil
	}
	if err := e.validate.Struct(params); err != nil {
		return fmt.Errorf("invalid request params: %w", err)
	}
	return nil
}

func (e *Encoder) encodePath(uri string, params any) (string, error) {
	if params == nil {
		return uri, nil
	}
	val, err := e.pathEncoder.Encode(&params)
	if err != nil {
		return "", fmt.Errorf("error encoding path params: %w", err)
	}

	pathParams := map[string]string{}
	for k, v := range val {
		pathParams[k] = v[0] // only accept the first one for a given key
	}

	for k, v := range pathParams {
		uri = strings.ReplaceAll(uri, fmt.Sprintf(":%s", k), url.PathEscape(v))
	}

	return uri, nil
}

func (e *Encoder) encodeQuery(params any) (string, error) {
	if params == nil {
		return "", nil
	}
	query, err := e.queryEncoder.Encode(&params)
	if err != nil {
		return "", fmt.Errorf("error encoding query params: %w", err)
	}
	return query.Encode(), nil
}

func newEncoder(tagName string) *form.Encoder {
	e := form.NewEncoder()
	e.SetMode(form.ModeExplicit)
	e.SetTagName(tagName)

	return e
}

func isDay(t time.Time) bool {
	if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 || t.Nanosecond() != 0 {
		return false
	}
	return true
}
