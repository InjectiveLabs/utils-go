package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bugsnag/bugsnag-go/errors"
	log "github.com/xlab/suplog"
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

type TokenMetadataResponse struct {
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

func getTokenMetaFromAlchemyByAddress(ctx context.Context, address string) *Meta {
	if address == "" {
		orPanicf(errors.Errorf("address is empty"), "invalid address")
	}

	mainnetAddressHex := KovanAddressMap[address]

	if mainnetAddressHex != "" {
		log.Infof("fetching token meta from Alchemy using [%s] instead of [%s]\n", mainnetAddressHex, address)
		address = mainnetAddressHex
	}

	reqBody := &AlchemyRequest{
		JsonRPC: "2.0",
		Method:  "alchemy_getTokenMetadata",
		ID:      randID(),
		Params: []string{
			address,
		},
	}

	reqData, err := json.Marshal(reqBody)
	orPanicf(err, "failed to marshal request body, body: [%+v]\n", reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(alchemyEndpoint, alchemyAPIKey), bytes.NewReader(reqData))
	orPanicf(err, "failed to create new request")

	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
	resp, err := cli.Do(req)
	orPanicf(err, "failed to send request")

	respData, err := io.ReadAll(io.LimitReader(resp.Body, 1000*1000))
	orPanicf(err, "failed to read data from resp.body")
	_ = resp.Body.Close()

	if resp.StatusCode != 200 {
		orPanicf(fmt.Errorf("error %d: %s", resp.StatusCode, respData), "")
	}

	var respBody *TokenMetadataResponse
	err = json.Unmarshal(respData, &respBody)
	orPanicf(err, "failed to unmarshal resp data")

	if respBody.Error != nil && respBody.Error.Code != 0 {
		if respBody.Error.Code == -32602 {
			panic(fmt.Sprintf("AlchemyAPI: contract address not found, address: [%s]", address))
		}
		panic(fmt.Sprintf("error %d: %s", respBody.Error.Code, respBody.Error.Message))
	}
	return respBody.Result
}
