package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/InjectiveLabs/sdk-go/client/common"
	"github.com/InjectiveLabs/sdk-go/client/exchange"
	injective_derivative_exchange_rpcpb "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
	injective_spot_exchange_rpcpb "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	"github.com/InjectiveLabs/utils-go/pkg/token"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type TokenFile map[string]*token.Token

func main() {
	tokenFile := flag.String("token-file", "token.json", "path to token file")
	marketSkipList := flag.String("market-skip-list", "", "comma separated list of markets to skip")
	flag.Parse()

	skipMarkets := map[string]struct{}{}
	for _, m := range strings.Split(*marketSkipList, ",") {
		skipMarkets[m] = struct{}{}
	}
	// read token file
	tokens, err := os.ReadFile(*tokenFile)
	if err != nil {
		log.Fatalf("cannot read token file: %s", err)
	}
	isValid := json.Valid(tokens)
	if !isValid {
		log.Fatalf("invalid json: %s", tokens)
	}

	var metadata = TokenFile{}

	decoder := json.NewDecoder(strings.NewReader(string(tokens)))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&TokenFile{})
	if err != nil {
		log.Fatalf("cannot decode token file: %s", err)
	}

	if err := json.Unmarshal(tokens, &metadata); err != nil {
		log.Fatalf("cannot unmarshal token file: %s", err)
	}

	log.Infof("token file is semantically valid")

	networkNode := "mainnet,lb"

	networkNodeSplit := strings.Split(networkNode, ",")
	networkStr, node := networkNodeSplit[0], networkNodeSplit[1]
	network := common.LoadNetwork(networkStr, node)

	exchangeClient, err := exchange.NewExchangeClient(network)
	if err != nil {
		log.Errorf("cannot create mainnet exchange client: %s", err)
		os.Exit(1)
	}

	networkStr = "testnet"
	network = common.LoadNetwork(networkStr, node)
	testnetExchangeClient, err := exchange.NewExchangeClient(network)
	if err != nil {
		log.Errorf("cannot create testnet exchange client: %s", err)
		os.Exit(1)
	}

	finder := MarketFinder{
		exchangeClient: map[string]exchange.ExchangeClient{
			"mainnet": exchangeClient,
			"testnet": testnetExchangeClient,
		},
	}
	for key, m := range metadata {
		if strings.Contains(strings.ToLower(key), "test") {
			continue
		}
		if m.Denom != "" {
			for _, market := range m.InjectiveMarkets {
				if _, ok := skipMarkets[market]; ok {
					continue
				}
				_, err = finder.findMarketOnAllNetwork(context.TODO(), market)
				if err != nil {
					log.Fatalf("%s: market %s not found: %s. If this is expected is possible to whitelist this market using `--market-skip-list` flag.", key, market, err)
				}
				log.Infof("%s: market %s found", key, market)
			}
		}
	}
	log.Infof("token file content is valid")
}

type MarketFinder struct {
	exchangeClient map[string]exchange.ExchangeClient
}

func (f MarketFinder) findMarketOnAllNetwork(ctx context.Context, market string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	var err error
	for key, client := range f.exchangeClient {
		spotMarketInfo := injective_spot_exchange_rpcpb.MarketResponse{}
		derivativeMarketInfo := injective_derivative_exchange_rpcpb.MarketResponse{}
		spotMarketInfo, err = client.GetSpotMarket(ctx, market)
		if err != nil {
			derivativeMarketInfo, err = client.GetDerivativeMarket(ctx, market)
			if err != nil &&
				!strings.Contains(err.Error(), "not found") &&
				!strings.Contains(err.Error(), "502") {
				return "", fmt.Errorf("cannot get market %s on %s: %s", market, key, err)
			}
			if err != nil {
				continue
			}
			return derivativeMarketInfo.Market.Ticker, nil
		}
		return spotMarketInfo.Market.Ticker, nil
	}
	if err != nil {
		return "", fmt.Errorf("market %s not found: %s", market, err)
	}
	return "", fmt.Errorf("market %s not found", market)
}
