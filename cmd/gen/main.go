package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/bugsnag/bugsnag-go/errors"
	log "github.com/InjectiveLabs/suplog"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	readEnv()
	alchemyAPIKey = os.Getenv(alchemyAPIKeyEnvVar)
	coinmarketcapAPIKey = os.Getenv(coinmarketcapAPIKeyEnvVar)
}

func main() {
	ctx := context.Background()

	err := os.RemoveAll(tokenMetaFileOutput)
	orPanicf(err, "failed to remove the output meta file")

	tokenMetaFileSource, err := os.OpenFile(tokenMetaFileSource, os.O_RDWR, os.ModePerm)
	orPanicf(err, "failed to open support token list file\n")
	defer func() {
		orPanicf(tokenMetaFileSource.Close(), "failed to close support token list file\n")
	}()
	fileContent, err := ioutil.ReadAll(tokenMetaFileSource)
	orPanicf(err, "failed to read support token list file\n")


	tokenMetaMap := TokenMetaMap{}
	err = json.Unmarshal(fileContent, &tokenMetaMap)
	orPanicf(err, "failed to json unmarshal token meta file content\n")

	log.Infof("got token meta map, [%d] tokens' metadata need to be filled\n", len(tokenMetaMap))
	log.Infof(logDivider)
	tokenMetaMap.tidy()

	// fill metas for each
	for token := range tokenMetaMap {
		err := fillTokenMeta(ctx, token, tokenMetaMap[token])
		orPanicf(err, "failed to fill token meta for [%s]\n", token)
		log.Infof("successfully filled token meta for [%s]\n", token)
	}
	log.Infof("finished fetching tokens' metadata\n")
	log.Infof(logDivider)

	tokenMetaMap.check()

	// write token metadata map to file
	newFileContent, err := json.MarshalIndent(&tokenMetaMap, "", "  ")
	orPanicf(err, "failed to json marshal new token meta map\n")

	tokenMetaFileOutput, err := os.Create(tokenMetaFileOutput)
	orPanicf(err, "failed to create the output meta file")

	_, err = tokenMetaFileOutput.Write(newFileContent)
	orPanicf(err, "failed to write new file content\n")
	orPanicf(tokenMetaFileOutput.Sync(), "failed to sync new file\n")

	log.Infof("successfully gen token meta file\n")
}

func orPanicf(err error, format string, args ...interface{}) {
	if err != nil {
		log.WithError(err).Panicf(format, args...)
	}
}

func fillTokenMeta(ctx context.Context, symbol string, t *Token) error {
	if t == nil {
		return errors.Errorf("empty token\n")
	}

	// 1. check coinGeckoID, since we need that to query the token's price
	if t.CoingeckoID == "" {
		return errors.Errorf("empty coingecko id, might cause an error when query token's price\n")
	} else {
		coin := GetCoingeckoTokenDetail(t.CoingeckoID)
		if strings.ToLower(coin.Platforms[ethereum]) != strings.ToLower(t.Address) {
			log.Warningf("token [%s] address [%s] is not same as in coingecko resp [%s], platforms: [%+v]\n",
				symbol, t.Address, coin.Platforms[ethereum], coin.Platforms)
		}
	}

	switch t.MetaSource {
	case MetaSourceAlchemy:
		if t.Address == "" {
			log.Infof("token [%s] doesn't have erc address, query that from CoinMarketCap by symbol", symbol)
			addr := GetEthereumAddressFromCoinMarketCapBySymbol(symbol)
			if addr == "" {
				log.Panicf("cannot solve ethereum address from symbol [%s], better to cover this in customizers\n", symbol)
			}
			log.Infof("got erc address [%s] for token [%s]", addr, symbol)
			t.Address = strings.ToLower(addr)
		}

		metadata := getTokenMetaFromAlchemyByAddress(ctx, t.Address)
		if metadata == nil {
			log.Panicf("token metadata is empty, address: [%s]\n", t.Address)
		}
		// sometimes, alchemy doesn't have logo for this symbol, search the logo url from CoinMarketCap
		if metadata.Logo == "" {
			logo := GetLogoBySymbol(symbol)
			if logo == "" {
				log.Warningf("cannot find logo for symbol [%s] from either Alchemy or CoinMarketCap", symbol)
			}
			metadata.Logo = logo
		}
		t.Meta = metadata
	case MetaSourceCustom:
		// this means the token's meta is already set in the json file manually, just validate that
		if t.Meta == nil {
			log.Warningf("token [%s] meta source is custom while the meta is empty")
		}
		if t.Meta.Name == "" {
			log.Warningf("token [%s] meta source is custom while the name is empty")
		}
		if t.Meta.Symbol == "" {
			log.Warningf("token [%s] meta source is custom while the symbol is empty")
		}
		if t.Meta.Decimals == 0 {
			log.Warningf("token [%s] meta source is custom while the decimal is zero")
		}
		if t.Meta.Logo == "" {
			log.Warningf("token [%s] meta source is custom while the logo is empty")
		}
	default:
		log.Panicf("unknown meta source [%s] for symbol [%s]", t.MetaSource, symbol)
	}

	return nil
}

func readEnv() {
	if envdata, _ := ioutil.ReadFile(".env"); len(envdata) > 0 {
		s := bufio.NewScanner(bytes.NewReader(envdata))
		for s.Scan() {
			txt := s.Text()
			valIdx := strings.IndexByte(txt, '=')
			if valIdx < 0 {
				continue
			}

			strValue := strings.Trim(txt[valIdx+1:], `"`)
			if err := os.Setenv(txt[:valIdx], strValue); err != nil {
				log.WithField("name", txt[:valIdx]).WithError(err).Warningln("failed to override ENV variable")
			}
		}
	}
}
