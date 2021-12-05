package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const coinInfoURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/info"
const coinmarketcapAPIKeyEnvVar = "COINMARKETCAP_API_KEY"
var coinmarketcapAPIKey string

func GetCoinInfoFromCoinMarketCap(symbol string) *Data {
	client := &http.Client{}
	req, err := http.NewRequest("GET", coinInfoURL, nil)
	orPanicf(err, "failed to new request to get coin info from CoinMarketCap\n")

	q := url.Values{}
	q.Add("symbol", symbol)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", coinmarketcapAPIKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	orPanicf(err, "failed to send request to CoinMarketCap\n")
	respBody, err := ioutil.ReadAll(resp.Body)
	orPanicf(err, "failed to read body from coin market cap resp")
	coinInfoResp := &CoinInfoResp{}
	err = json.Unmarshal(respBody, coinInfoResp)
	orPanicf(err, "failed to json unmarshal resp body")
	for _, data := range coinInfoResp.Data {
		return data
	}
	return nil
}

const PlatformEthereum = "Ethereum"

func GetEthereumAddressFromCoinMarketCapBySymbol(symbol string) string {
	data := GetCoinInfoFromCoinMarketCap(symbol)
	if data == nil || len(data.ContractAddress) == 0 {
		return ""
	}
	for _, ca := range data.ContractAddress {
		if ca == nil {
			continue
		}
		if ca.Platform.Name == PlatformEthereum {
			return ca.ContractAddress
		}
	}
	return ""
}

func GetLogoBySymbol(symbol string)string{
	data := GetCoinInfoFromCoinMarketCap(symbol)
	if data == nil || len(data.ContractAddress) == 0 {
		return ""
	}
	return data.Logo
}

type CoinInfoResp struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
		Notice       interface{} `json:"notice"`
	} `json:"status"`
	Data map[string]*Data `json:"data"`
}

type Data struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Logo        string   `json:"logo"`
	Subreddit   string   `json:"subreddit"`
	Notice      string   `json:"notice"`
	Tags        []string `json:"tags"`
	TagNames    []string `json:"tag-names"`
	TagGroups   []string `json:"tag-groups"`
	Urls        struct {
		Website      []string      `json:"website"`
		Twitter      []string      `json:"twitter"`
		MessageBoard []string      `json:"message_board"`
		Chat         []string      `json:"chat"`
		Facebook     []interface{} `json:"facebook"`
		Explorer     []string      `json:"explorer"`
		Reddit       []string      `json:"reddit"`
		TechnicalDoc []string      `json:"technical_doc"`
		SourceCode   []string      `json:"source_code"`
		Announcement []interface{} `json:"announcement"`
	} `json:"urls"`
	Platform                      interface{}        `json:"platform"`
	DateAdded                     time.Time          `json:"date_added"`
	TwitterUsername               string             `json:"twitter_username"`
	IsHidden                      int                `json:"is_hidden"`
	DateLaunched                  interface{}        `json:"date_launched"`
	ContractAddress               []*ContractAddress `json:"contract_address"`
	SelfReportedCirculatingSupply interface{}        `json:"self_reported_circulating_supply"`
	SelfReportedTags              []string           `json:"self_reported_tags"`
}
type ContractAddress struct {
	ContractAddress string `json:"contract_address"`
	Platform        struct {
		Name string `json:"name"`
		Coin struct {
			Id     string `json:"id"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
			Slug   string `json:"slug"`
		} `json:"coin"`
	} `json:"platform"`
}
