package token

// TODO: this file is a copy from submodule, bad practice totally

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type AlchemyRequest struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`

	Method string   `json:"method"`
	Params []string `json:"params"`
}

type MetadataResponse struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`

	Error  *RPCError `json:"error"`
	Result *Meta     `json:"result"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func randID() int {
	return int(1 + rand.Int63n(100000000))
}

func getTokenMetaFromAlchemyByAddress(ctx context.Context, address string) (*Meta, error) {
	reqBody := &AlchemyRequest{
		JsonRPC: "2.0",
		Method:  "alchemy_getTokenMetadata",
		ID:      randID(),
		Params: []string{
			address,
		},
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal request body, body: [%+v]\n", reqBody)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", alchemyEndpoint, bytes.NewReader(reqData))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new request")
	}

	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send request")
	}

	respData, err := io.ReadAll(io.LimitReader(resp.Body, 1000*1000))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read data from resp.body")
	}
	_ = resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("error %d: %s", resp.StatusCode, respData)
	}

	var respBody *MetadataResponse
	err = json.Unmarshal(respData, &respBody)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal resp data")
	}

	if respBody.Error != nil && respBody.Error.Code != 0 {
		if respBody.Error.Code == -32602 {
			return nil, errors.Errorf("AlchemyAPI: contract address not found, address: [%s]", address)
		}
		return nil, errors.Errorf("error %d: %s", respBody.Error.Code, respBody.Error.Message)
	}
	return respBody.Result, nil
}
