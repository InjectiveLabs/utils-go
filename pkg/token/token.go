package token

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/xlab/suplog"
	"strings"
	"time"
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
	symbolMap = make(map[string]*Token, len(tokenMap))
	for s := range tokenMap {
		symbolMap[strings.ToLower(s)] = tokenMap[s]
	}
	// for no case sensitivity, and addresses in json file have no prefix "peggy"
	addressMap = make(map[string]*Token, len(tokenMap))
	for s := range tokenMap {
		addressMap[strings.ToLower(tokenMap[s].Address)] = tokenMap[s]
	}
	log.Infof("successfully loaded token meta config\n")
}

// GetTokenBySymbol no case sensitivity, USD/usd/Usd are all fine
func GetTokenBySymbol(symbol string) *Token {
	return symbolMap[strings.ToLower(symbol)]
}

// GetTokenByAddress no case sensitivity, and it's safe to pass address with prefix "peggy"
// for unknown address, request metadata from alchemy
func GetTokenByAddress(address string) *Token {
	address = strings.ToLower(strings.TrimPrefix(address, "peggy"))
	token, ok := addressMap[address]
	if ok && token != nil {
		return token
	}
	// token not exist in address map, search from alchemy
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()
	tokenMeta, err := getTokenMetaFromAlchemyByAddress(ctx, address)
	if err == nil && tokenMeta != nil {
		return &Token{Address: address, Meta: tokenMeta}
	}
	return nil
}
