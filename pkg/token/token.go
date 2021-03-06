package token

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/InjectiveLabs/suplog"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed token_meta.json
var tokenMetaFileContent []byte

type EthereumAddress struct {
	common.Address
}

var symbolMapLock sync.RWMutex
var symbolMap map[string]*Token
var addressMapLock sync.RWMutex
var addressMap map[string]*Token
var denomMapLock sync.RWMutex
var denomMap map[string]*Token

const alchemyEndpoint = "https://eth-mainnet.alchemyapi.io/v2/%s"
const alchemyAPIKeyEnvVar = "ALCHEMY_API_KEY"

var alchemyAPIKey string

func init() {
	alchemyAPIKey = os.Getenv(alchemyAPIKeyEnvVar)
	var tokenMap Dict
	err := json.Unmarshal(tokenMetaFileContent, &tokenMap)
	if err != nil {
		panic(err)
	}
	// no case sensitivity
	symbolMapLock.Lock()
	addressMapLock.Lock()
	denomMapLock.Lock()
	symbolMap = make(map[string]*Token, len(tokenMap))
	addressMap = make(map[string]*Token, len(tokenMap))
	denomMap = make(map[string]*Token, len(tokenMap))
	for s, token := range tokenMap {
		if token == nil || s == "" {
			log.Warningf("got invalid token metadata, symbol: [%s]", s)
			continue
		}
		symbolMap[strings.ToLower(s)] = token
		addressMap[strings.ToLower(token.Address)] = token
		denomMap[strings.ToLower(token.Denom)] = token
	}
	symbolMapLock.Unlock()
	addressMapLock.Unlock()
	denomMapLock.Unlock()
	log.Infof("successfully loaded token meta config\n")
}

// GetTokenBySymbol no case sensitivity, USD/usd/Usd are all fine
func GetTokenBySymbol(symbol string) *Token {
	symbolMapLock.RLock()
	defer symbolMapLock.RUnlock()
	return symbolMap[strings.ToLower(symbol)]
}

// GetTokenByAddress no case sensitivity, and it's safe to pass address with prefix "peggy"
// for unknown address, request metadata from alchemy
func GetTokenByAddress(address string) *Token {
	if strings.ToLower(address) == "inj" {
		return GetTokenBySymbol("inj")
	}
	address = strings.ToLower(strings.TrimPrefix(address, "peggy"))
	addressMapLock.RLock()
	token, ok := addressMap[address]
	addressMapLock.RUnlock()
	if ok && token != nil {
		return token
	}
	// token not exist in address map, search from alchemy
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()
	tokenMeta, err := getTokenMetaFromAlchemyByAddress(ctx, address)
	if err == nil && tokenMeta != nil {
		// add token metadata into the cache map
		addressMapLock.Lock()
		addressMap[address] = &Token{Address: address, Meta: tokenMeta}
		addressMapLock.Unlock()
		return &Token{Address: address, Meta: tokenMeta}
	}
	return nil
}

// GetTokenByDenom no case sensitivity
func GetTokenByDenom(denom string) *Token {
	denomMapLock.RLock()
	defer denomMapLock.RUnlock()
	return denomMap[strings.ToLower(denom)]
}
