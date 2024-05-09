package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"go.tradeforge.dev/oxr/internal/encoder"
	"go.tradeforge.dev/oxr/model"
	"go.tradeforge.dev/oxr/util"

	"github.com/go-resty/resty/v2"
)

const clientVersion = "v0.0.0"

const (
	DefaultRetryCount    = 3
	DefaultClientTimeout = 10 * time.Second
)

type EventReader interface {
	Listen(ctx context.Context, stream io.Reader) (<-chan Event, <-chan error)
}

type Event interface {
	GetData() []byte
	GetTimestamp() time.Time
}

// Client defines an HTTP client for the Polygon REST API.
type Client struct {
	HTTP    *resty.Client
	encoder *encoder.Encoder

	logger *slog.Logger
}

// New returns a new client with the specified API key and config.
func New(
	apiURL string,
	logger *slog.Logger,
) *Client {
	if logger == nil {
		logger = slog.New(&util.NilHandler{})
	}
	return newClient(apiURL, logger)
}

func newClient(
	apiURL string,
	logger *slog.Logger,
) *Client {
	c := resty.New()

	c.SetBaseURL(apiURL)
	c.SetRetryCount(DefaultRetryCount)
	c.SetHeader("User-Agent", fmt.Sprintf("Open Exchange Rates client/%v", clientVersion))
	c.SetHeader("Accept", "application/json")

	return &Client{
		HTTP:    c,
		encoder: encoder.New(),
		logger:  logger,
	}
}

func (c *Client) SetAuthScheme(scheme string) *Client {
	c.HTTP.SetAuthScheme(scheme)
	return c
}

func (c *Client) SetAuthToken(apiKey string) *Client {
	c.HTTP.SetAuthToken(apiKey)
	return c
}

func (c *Client) SetHeader(key, value string) *Client {
	c.HTTP.SetHeader(key, value)
	return c
}

func (c *Client) SetLogger(logger *slog.Logger) *Client {
	c.logger = logger
	return c
}

// Call makes an API call based on the request params and options. The response is automatically unmarshalled.
func (c *Client) Call(ctx context.Context, method, path string, params, response any, opts ...model.RequestOption) error {
	uri, err := c.encoder.EncodeParams(path, params)
	if err != nil {
		return err
	}
	return c.CallURL(ctx, method, uri, response, opts...)
}

// CallURL makes an API call based on a request URI and options. The response is automatically unmarshalled.
func (c *Client) CallURL(ctx context.Context, method, uri string, response any, opts ...model.RequestOption) error {
	options := mergeOptions(opts...)

	c.HTTP.SetTimeout(DefaultClientTimeout)
	req := c.HTTP.R().SetContext(ctx)
	if options.Body != nil {
		b, err := json.Marshal(options.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %w", err)
		}
		req.SetBody(b)
	}
	req.SetQueryParamsFromValues(options.QueryParams)
	req.SetHeaderMultiValues(options.Headers)
	req.SetResult(response).SetError(&model.ResponseError{})
	req.SetHeader("Content-Type", "application/json")

	res, err := req.Execute(method, uri)
	if err != nil {
		c.logger.Error(
			err.Error(),
			slog.Any("response", res))
		return fmt.Errorf("failed to execute request: %w", err)
	} else if res.IsError() {
		c.logger.Error(
			res.String(),
			slog.Any("response", res))
		responseError := parseResponseError(res)
		return responseError
	}

	if options.Trace {
		sanitizedHeaders := req.Header
		for k := range sanitizedHeaders {
			if k == "Authorization" {
				sanitizedHeaders[k] = []string{"REDACTED"}
			}
		}
		c.logger.Debug(
			"request",
			slog.String("url", uri),
			slog.Any("request headers", sanitizedHeaders),
			slog.Any("response headers", res.Header()),
		)
	}

	return nil
}

func mergeOptions(opts ...model.RequestOption) *model.RequestOptions {
	options := &model.RequestOptions{}
	for _, o := range opts {
		o(options)
	}

	return options
}

func parseResponseError(res *resty.Response) *model.ResponseError {
	if res == nil {
		return nil
	}
	var responseError *model.ResponseError
	if !errors.As(res.Error().(error), &responseError) {
		panic("type assertion of res.Error() failed: res.Error() does not implement error or *model.ResponseError")
	}
	responseError.RequestID = res.Header().Get("X-Request-ID")
	responseError.StatusCode = res.StatusCode()
	b := struct {
		Message string `json:"message"`
	}{}
	if err := json.Unmarshal(res.Body(), &b); err != nil {
		responseError.Message = string(res.Body())
	} else {
		responseError.Message = b.Message
	}
	return responseError
}
