package token

import (
	_ "embed"
	"encoding/json"
	log "github.com/xlab/suplog"
	"strings"
)

// this path has to be hardcoded, no other ways
//go:generate cp ../../lib/token-meta/meta/token_meta.json ./token_meta.json
//go:embed token_meta.json
var tokenMetaFileContent []byte

type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Address     string `json:"address"`
	CoinGeckoId string `json:"coinGeckoId"`
}

var symbolMap map[string]*Token
var addressMap map[string]*Token

func init() {
	err := json.Unmarshal(tokenMetaFileContent, &symbolMap)
	if err != nil {
		panic(err)
	}
	// for no case sensitivity
	for s := range symbolMap {
		symbolMap[strings.ToLower(s)] = symbolMap[s]
	}
	addressMap = map[string]*Token{}
	for s := range symbolMap {
		addressMap[symbolMap[s].Address] = symbolMap[s]
	}
	log.Infof("successfully loaded token meta config\n")
}

// GetTokenMetaBySymbol no case sensitivity, USD/usd/Usd are all fine
func GetTokenMetaBySymbol(symbol string) *Token {
	return symbolMap[strings.ToLower(symbol)]
}

func GetTokenMetaByAddress(address string) *Token {
	if strings.HasPrefix(address, "peggy") {
		address = address[5:]
	}
	return addressMap[address]
}
