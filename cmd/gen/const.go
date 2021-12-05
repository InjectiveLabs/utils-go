package main

const tokenMetaFileSource = "pkg/token/token_meta.json"
const tokenMetaFileOutput = "cmd/gen/token_meta.json"

const alchemyEndpoint = "https://eth-mainnet.alchemyapi.io/v2/%s"
const alchemyAPIKeyEnvVar = "ALCHEMY_API_KEY"
var alchemyAPIKey string
const (
	KovanDAIAddress     = "0x9566902A13ce8aD8c730743e54ca0fF3657470a0"
	MainnetDAIAddress   = "0x6b175474e89094c44da98b954eedeac495271d0f"
	KovanWETHAddress    = "0xd0A1E359811322d97991E03f863a0C30C2cF029C"
	MainnetWETHAddress  = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	KovanSUSHIAddress   = "0x12D04f84d49FE289B2BD65404dB6504789BE34B1"
	MainnetSUSHIAddress = "0x6b3595068778dd592e39a122f4f5a5cf09c90fe2"
	KovanWBTCAddress    = "0x0568b026599c2FD48AA1D80A16122db2F9c7Ae15"
	MainnetWBTCAddress  = "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"
	KovanSNXAddress     = "0x7FDb81B0b8a010dd4FFc57C3fecbf145BA8Bd947"
	MainnetSNXAddress   = "0xc011a73ee8576fb46f5e1c5751ca3b9fe0af2a6f"
	KovanUSDTAddress    = "0x69efCB62D98f4a6ff5a0b0CFaa4AAbB122e85e08"
	MainnetUSDTAddress  = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	KovanUSDCAddress    = "0xc83DCEA3Ec44b7D3Ec70690BAb1e6292A80e6DC3"
	MainnetUSDCAddress  = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	KovanUNIAddress     = "0x138b989687da853a561D4edE88D8281434211780"
	MainnetUNIAddress   = "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"
	KovanAAVEAddress    = "0x69BeD9289Eb970F021BA86fec646f9C427E0320A"
	MainnetAAVEAddress  = "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9"
	KovanMATICAddress   = "0x724d7E46BF2cC15de3932F547a60018C286312A7"
	MainnetMATICAddress = "0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0"
	KovanZRXAddress     = "0xB4EF9D74108980FecE40d9205c3d1c94090a3b50"
	MainnetZRXAddress   = "0xe41d2489571d322189246dafa5ebde1f4699f498"
	KovanLINKAddress    = "0xC843F43093f8d32c01a065ed2a0a34fb54BAaf3F"
	MainnetLINKAddress  = "0x514910771af9ca656af840dff83e8264ecf986ca"
	KovanBNBAddress     = "0xf833cAd2b46B49Ef96244B974AAFf8B80Ff84FDd"
	MainnetBNBAddress   = "0xB8c77482e45F1F44dE1745F52C74426C631bDD52"
	KovanYFIAddress     = "0x6aCd36eB845a8f905512D5F259c1233242349266"
	MainnetYFIAddress   = "0x0bc529c00c6401aef6d220be8c6ea1667f6ad93e"
)

var KovanAddressMap = map[string]string{
	KovanDAIAddress:   MainnetDAIAddress,
	KovanWETHAddress:  MainnetWETHAddress,
	KovanSUSHIAddress: MainnetSUSHIAddress,
	KovanWBTCAddress:  MainnetWBTCAddress,
	KovanSNXAddress:   MainnetSNXAddress,
	KovanUSDTAddress:  MainnetUSDTAddress,
	KovanUSDCAddress:  MainnetUSDCAddress,
	KovanUNIAddress:   MainnetUNIAddress,
	KovanAAVEAddress:  MainnetAAVEAddress,
	KovanMATICAddress: MainnetMATICAddress,
	KovanZRXAddress:   MainnetZRXAddress,
	KovanLINKAddress:  MainnetLINKAddress,
	KovanBNBAddress:   MainnetBNBAddress,
	KovanYFIAddress:   MainnetYFIAddress,
}

const ethereum = "ethereum"

const logDivider = "============================================================================\n"