package constants

const (
	NETWORK            string = "preprod"
	BFC_NETWORK_ID     int    = 0
	MAESTRO_NETWORK_ID int    = 3 // preprod
	BFC_API_URL        string = "https://cardano-preprod.blockfrost.io/api"
	BFC_API_KEY        string = "preprod9zzl4g8Xa3faU50a1OVDZdPeQ92ZsdcT"
	MAESTRO_API_KEY    string = "Dq21Cy4YQeH7AL8q61wAqYcHHur7QT9S"
	OGMIGO_ENDPOINT    string = "ws://localhost:1337"
	KUGO_ENDPOINT      string = "http://localhost:1442"

	EKIVAL_FEE int = 25_000_000

	UNCOMMITTED_ORDER_STATUS string = "UNCOMMITTED_ORDER"
	COMMITTED_ORDER_STATUS   string = "COMMITTED_ORDER"
	REMIT_CONFIRMED_STATUS   string = "REMIT_CONFIRMED"

	ADA_P2P_BUY_FEE_TYPE    int = 1
	ADA_P2P_SELL_FEE_TYPE   int = 2
	TOKEN_P2P_BUY_FEE_TYPE  int = 3
	TOKEN_P2P_SELL_FEE_TYPE int = 4
	STAKING_FEE_TYPE        int = 5
	BUY_OFFER_FEE_TYPE      int = 6
	SELL_OFFER_FEE_TYPE     int = 7

	INDEX_ZERO  uint64 = 121
	INDEX_ONE   uint64 = 122
	INDEX_TWO   uint64 = 123
	INDEX_THREE uint64 = 124
	INDEX_FOUR  uint64 = 125
	INDEX_FIVE  uint64 = 126
	INDEX_SIX   uint64 = 127
	INDEX_SEVEN uint64 = 1280
	INDEX_EIGHT uint64 = 1281
	INDEX_NINE  uint64 = 1282
	INDEX_TEN   uint64 = 1283
)
