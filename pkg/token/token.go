package token

import (
	"context"
	_ "embed"
	"encoding/json"

	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	log "github.com/InjectiveLabs/suplog"
	"github.com/ethereum/go-ethereum/common"
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

const cacheTTL = time.Hour * 3
const cacheRefreshInterval = time.Minute * 10

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
	if alchemyAPIKey != "" {
		cacheCleaner()
	}
}

// GetTokenBySymbol no case sensitivity, USD/usd/Usd are all fine
func GetTokenBySymbol(symbol string) *Token {
	symbolMapLock.RLock()
	defer symbolMapLock.RUnlock()
	return symbolMap[strings.ToLower(symbol)]
}

// GetTokenByAddress no case sensitivity, and it's safe to pass address with prefix "peggy"
// for unknown address, request metadata from alchemy
// This method rely on an internal cache, so it's safe to call it frequently
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

	return nil
}

// GetTokenByDenom no case sensitivity
func GetTokenByDenom(denom string) *Token {
	denomMapLock.RLock()
	defer denomMapLock.RUnlock()
	return denomMap[strings.ToLower(denom)]
}

func SetDenomToMetadataMap(metadataMap map[string]*Token) {
	denomMapLock.Lock()
	defer denomMapLock.Unlock()
	denomMap = metadataMap
}

func cacheCleaner() {
	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		defer stop()
		ticker := time.NewTicker(cacheRefreshInterval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				addressMapLock.Lock()
				for k, v := range addressMap {
					if time.Since(v.LastAccessTime) > cacheTTL {
						delete(addressMap, k)
					}
				}
				addressMapLock.Unlock()
			}
		}
	}()
}
