package jsonrpc

import (
	"bytes"
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/go-faster/errors"
)

type Request = http.Request

func Prepare[Params any](ctx context.Context, url, method string, params Params) (*Request, error) {
	buff := bytes.NewBuffer([]byte{})

	encoder := sonic.ConfigDefault.NewEncoder(buff)
	if err := encoder.Encode(createRequest[Params](method, params)); err != nil {
		return nil, errors.Wrap(err, "encode request")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buff)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	return request, nil
}

func Execute[Result any](request *Request) (*Result, error) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "execute request")
	}

	decoder := sonic.ConfigDefault.NewDecoder(response.Body)
	defer func() {
		_ = response.Body.Close()
	}()

	var result RPCResponse[Result]

	if err := decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	return &result.Result, nil
}
