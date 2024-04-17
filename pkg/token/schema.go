package token

import "time"

type MetaSource string

const (
	MetaSourceAlchemy MetaSource = "alchemy" // query from Alchemy API
	MetaSourceCustom  MetaSource = "custom"  // defined in Customizers
)

type Token struct {
	Address          string     `json:"address"`
	CoingeckoID      string     `json:"coinGeckoId"`
	Denom            string     `json:"denom,omitempty"`
	MetaSource       MetaSource `json:"metaSource"`
	Meta             *Meta      `json:"meta"`
	LastAccessTime   time.Time  `json:"-"`
	InjectiveMarkets []string   `json:"injectiveMarkets"`
	StartPrice       float64    `json:"start_price"`
}

// Meta this struct is the same as the metadata in the resp of Alchemy
type Meta struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
}

type Dict map[string]*Token
