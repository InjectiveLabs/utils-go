package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CoinGeckoGetByIdResp struct {
	Id              string            `json:"id"`
	Symbol          string            `json:"symbol"`
	Name            string            `json:"name"`
	AssetPlatformId string            `json:"asset_platform_id"`
	Platforms       map[string]string `json:"platforms"`
}

func GetCoingeckoTokenDetail(id string) *CoinGeckoGetByIdResp {
	resp, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers&false&market_data=false&community_data=false&developer_data&false&sparkline=false", id))
	orPanicf(err, "failed get coin info from coingecko with id [%s]", id)
	body, err := ioutil.ReadAll(resp.Body)
	orPanicf(err, "failed to read all from coingecko resp body")
	token := new(CoinGeckoGetByIdResp)
	err = json.Unmarshal(body, token)
	orPanicf(err, "failed to json unmarshal coingecko resp body, body: [%s]\n", string(body))
	time.Sleep(1200 * time.Millisecond)
	return token
}
