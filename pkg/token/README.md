# Token Metadata for Go

## Usage
Import this repo directly into your own go repo

## Maintain
### To add a new token's metadata:
1. add a new `token_name: {address:token_address, coinGeckoId:token_coingecko_id, denom:token_denom, metaSource: [Alchemy/custom]}` kv pair in `token_meta.json`
2. run `make gen`

### To add some new util functions
1. develop your new util functions
2. don't forget to run `make build` to make sure no error in compiling at least