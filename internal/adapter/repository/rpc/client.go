package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"watcher/configs"
)

type Client struct {
	url    string
	client *http.Client
}

type requestBody struct {
	ID      int64  `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any
}

type responseBody struct {
	JsonRPC string         `json:"jsonrpc"`
	Result  interface{}    `json:"result,omitempty"`
	Error   *responseError `json:"error,omitempty"`
	ID      int            `json:"id"`
}

type responseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpError struct {
	Code int
	err  error
}

// Error function is provided to be used as error object.
func (e *HttpError) Error() string {
	return e.err.Error()
}

func NewRpcClient(cfg configs.Config) *Client {
	return &Client{
		url:    cfg.NodeUrl,
		client: &http.Client{},
	}
}

func (cli *Client) CallFor(ctx context.Context, dest any, method string, params ...any) error {
	id := time.Now().UnixMicro()

	body := requestBody{
		ID:      id,
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	enc, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("rpc call %v() on %v: %w", method, cli.url, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cli.url, bytes.NewReader(enc))
	if err != nil {
		return fmt.Errorf("rpc call %v() on %v: %w", method, cli.url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cli.client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var rpcResponse *responseBody

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()

	err = decoder.Decode(&rpcResponse)

	if err != nil {
		// if we have some http error, return it
		if resp.StatusCode >= 400 {
			return &HttpError{
				Code: resp.StatusCode,
				err: fmt.Errorf(
					"rpc call %v() on %v status code: %v. could not decode body to rpc response: %w",
					body.Method,
					req.URL.Redacted(),
					resp.StatusCode, err,
				),
			}
		}
		return fmt.Errorf(
			"rpc call %v() on %v status code: %v. could not decode body to rpc response: %w",
			body.Method,
			req.URL.Redacted(),
			resp.StatusCode,
			err,
		)
	}

	// response body empty
	if rpcResponse == nil {
		// if we have some http error, return it
		if resp.StatusCode >= 400 {
			return &HttpError{
				Code: resp.StatusCode,
				err:  fmt.Errorf("rpc call %v() on %v status code: %v. rpc response missing", body.Method, req.URL.Redacted(), resp.StatusCode),
			}
		}
		return fmt.Errorf(
			"rpc call %v() on %v status code: %v. rpc response missing",
			body.Method,
			req.URL.Redacted(),
			resp.StatusCode,
		)
	}

	// if we have a response body, but also a http error situation, return both
	if resp.StatusCode >= 400 {
		if rpcResponse.Error != nil {
			return &HttpError{
				Code: resp.StatusCode,
				err: fmt.Errorf(
					"rpc call %v() on %v status code: %v. rpc response error: %v",
					body.Method,
					req.URL.Redacted(),
					resp.StatusCode,
					rpcResponse.Error,
				),
			}
		}
		return &HttpError{
			Code: resp.StatusCode,
			err: fmt.Errorf(
				"rpc call %v() on %v status code: %v. no rpc error available",
				body.Method,
				req.URL.Redacted(),
				resp.StatusCode,
			),
		}
	}

	err = rpcResponse.GetObject(dest)
	if err != nil {
		return err
	}

	return nil
}

// GetInt converts the rpc response to an int64 and returns it.
//
// If result was not an integer an error is returned.
func (r *responseBody) GetInt() (int64, error) {
	val, ok := r.Result.(json.Number)
	if !ok {
		return 0, fmt.Errorf("could not parse int64 from %s", r.Result)
	}

	i, err := val.Int64()
	if err != nil {
		return 0, err
	}

	return i, nil
}

// GetFloat converts the rpc response to float64 and returns it.
//
// If result was not an float64 an error is returned.
func (r *responseBody) GetFloat() (float64, error) {
	val, ok := r.Result.(json.Number)
	if !ok {
		return 0, fmt.Errorf("could not parse float64 from %s", r.Result)
	}

	f, err := val.Float64()
	if err != nil {
		return 0, err
	}

	return f, nil
}

// GetBool converts the rpc response to a bool and returns it.
//
// If result was not a bool an error is returned.
func (r *responseBody) GetBool() (bool, error) {
	val, ok := r.Result.(bool)
	if !ok {
		return false, fmt.Errorf("could not parse bool from %s", r.Result)
	}

	return val, nil
}

// GetString converts the rpc response to a string and returns it.
//
// If result was not a string an error is returned.
func (r *responseBody) GetString() (string, error) {
	val, ok := r.Result.(string)
	if !ok {
		return "", fmt.Errorf("could not parse string from %s", r.Result)
	}

	return val, nil
}

// GetObject converts the rpc response to an arbitrary type.
//
// The function works as you would expect it from json.Unmarshal()
func (r *responseBody) GetObject(toType interface{}) error {
	js, err := json.Marshal(r.Result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, toType)
	if err != nil {
		return err
	}

	return nil
}
