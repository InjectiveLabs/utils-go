package token

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Address     string `json:"address"`
	CoinGeckoId string `json:"coinGeckoId"`
}

var symbolMap map[string]*Token
var addressMap map[string]*Token

const tokenMetaFilePath = "lib/token-meta/meta"
const tokenMetaFileName = "token_meta.json"

func init() {
	// load token meta from submodule
	file, err := os.Open(path.Join(tokenMetaFilePath, tokenMetaFileName))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) { _ = f.Close() }(file)
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	symbolMap = map[string]*Token{}
	err = json.Unmarshal(fileContent, &symbolMap)
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
}

// GetTokenMetaBySymbol no case sensitivity, USD/usd/Usd are all fine
func GetTokenMetaBySymbol(symbol string) *Token {
	return symbolMap[strings.ToLower(symbol)]
}

func GetTokenMetaByAddress(address string) *Token {
	return symbolMap[address]
}
