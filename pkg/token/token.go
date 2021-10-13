package token

import (
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/xlab/suplog"
	"strings"
)

// this path has to be hardcoded, no other ways
//go:generate cp ../../lib/token-meta/meta/token_meta.json ./token_meta.json
//go:embed token_meta.json
var tokenMetaFileContent []byte

type EthereumAddress struct {
	common.Address
}

var symbolMap map[string]*Token
var addressMap map[string]*Token

func init() {
	var tokenMap Dict
	err := json.Unmarshal(tokenMetaFileContent, &tokenMap)
	if err != nil {
		panic(err)
	}
	// for no case sensitivity
	for s := range tokenMap {
		symbolMap[strings.ToLower(s)] = tokenMap[s]
	}
	// addresses in json file have no prefix "peggy"
	addressMap = map[string]*Token{}
	for s := range tokenMap {
		addressMap[tokenMap[s].Address] = tokenMap[s]
	}
	log.Infof("successfully loaded token meta config\n")
}

// GetTokenBySymbol no case sensitivity, USD/usd/Usd are all fine
func GetTokenBySymbol(symbol string) *Token {
	return symbolMap[strings.ToLower(symbol)]
}

// GetTokenByAddress will trim prefix "peggy"
func GetTokenByAddress(address string) *Token {
	return addressMap[strings.TrimPrefix(address, "peggy")]
}
