package main

import (
	log "github.com/xlab/suplog"
	"strings"
)

type MetaSource string

const (
	MetaSourceAlchemy MetaSource = "Alchemy" // query from Alchemy API
	MetaSourceCustom  MetaSource = "custom"  // defined in Customizers
)

type Token struct {
	Address     string     `json:"address"`
	CoingeckoID string     `json:"coinGeckoId"`
	Denom       string     `json:"denom,omitempty"`
	MetaSource  MetaSource `json:"metaSource"`
	Meta        *Meta      `json:"meta"`
}

// Meta this struct is the same as the metadata in the resp of Alchemy
type Meta struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
}

type TokenMetaMap map[string]*Token

func (m *TokenMetaMap) tidy() {
	for _, t := range *m {
		t.Address = strings.ToLower(t.Address)
		t.Denom = strings.ToLower(t.Denom)
	}
}

func (m *TokenMetaMap) check() {
	for symbol, t := range *m {
		if symbol == "" {
			log.Warningf("token symbol is empty\n")
		}
		if t.CoingeckoID == "" {
			log.Warningf("token [%s] coingecko id is empty\n", symbol)
		}
		if t.Address == "" && t.Denom == "" {
			log.Warningf("token [%s] both address and denom are empty\n", symbol)
		}
		if t.Meta == nil {
			log.Warningf("token [%s] meta is empty\n", symbol)
		}

		if t.Meta.Symbol == "" {
			log.Warningf("token [%s] meta.symbol is empty\n", symbol)
		}
		if t.Meta.Name == "" {
			log.Warningf("token [%s] meta.name is empty\n", symbol)
		}
		if t.Meta.Decimals == 0 {
			log.Warningf("token [%s] meta.decimal is 0\n", symbol)
		}
		if t.Meta.Logo == "" {
			log.Warningf("token [%s] meta.logo is empty\n", symbol)
		}
	}
}
