package tests

const (
	MAKER_ADDRESS    string = "addr_test1qqt85kcauy3uktlfldmhqn4dnn5vxgv8s97t729akje92f3z5uathwrlk6dwzy9j89lmsy2qp7evugtxmhf9pycw44rqj76x7l"
	TAKER_ADDRESS    string = "addr_test1qr43kyyys0sg0d8khhjq4y348zuc4mnzln836hxfrqjasx4rqghlrj99l5vrdmyrtg6mhkyxa88kwq5yf225a4m9pkes2a4xz3"
	DEBUGGER_ADDRESS string = "addr_test1qqd9lnkkk4uj96ftyqcgune8jq7g04wccmzcd3hehrddl8at3u4s70lrcaw39t522nq3y9953k2w47t4mgtj678xpp7q0dx95k"
	ADMIN_ADDRESS    string = "addr_test1qp0ccf96g40e45wtc98g4nvdxrgxwdwy0c4ymlpq9gprts5eqx40n59saqk2jlwc49uw30mfmlm7nckja2sg6wlrclmqaymgpe"

	REMIT_CONFIRMED_STATUS string = "REMIT_CONFIRMED"

	TEKI_POLICY_ID  string = "fe691a72d5a591a2ec8602a9904b30c005138d20af67b36087b27221"
	TEKI_TOKEN_NAME string = "tEKI"

	FEE_TREASURY_ADDRESS string = "addr_test1xz7zrmrn0dkf9utk7w4f0hecrns7w2lj85ce7fhnxj24awz2m3nsefutuywd8uqzd2n4u7w0vmymjhmmuuc46hdxy85qpmc4g3"

	ADA_P2P_BUY_ESCROW_ADDRESS                    string = "addr_test1xrgysjt2g0t4l7h54px984wehx9uhtautarly2nqq83w5622m3nsefutuywd8uqzd2n4u7w0vmymjhmmuuc46hdxy85qzlh2dc"
	ADA_P2P_BUY_ESCROW_SCRIPT_REF_UTXO_TXID       string = "87ce9218b91f09156b15bdb6fa5c3a8aa78d846184f2e1b5375df0114abd8b6d"
	ADA_P2P_BUY_ESCROW_SCRIPT_REF_UTXO_TXID_INDEX int    = 0
	APBST_POLICY_ID                               string = "aefdb5f954ea897ec536de425e62632728fa78b4638882654d0f4075"
	APBST_SCRIPT_REF_UTXO_TXID                    string = "fbf2cc43233490bc7f748957711bda85a530484113c8b65381b0613c30c13ac2"
	APBST_SCRIPT_REF_UTXO_TXID_INDEX              int    = 0

	ADA_P2P_BUY_FEE_TYPE int = 1

	NETWORK            string = "preprod"
	BFC_NETWORK_ID     int    = 0
	MAESTRO_NETWORK_ID int    = 3 // preprod
	BFC_API_URL        string = "https://cardano-preprod.blockfrost.io/api"
	BFC_API_KEY        string = "preprodLRHR4UwMuz7TuVOGtRJPcUIvcBDjy9Oz"
	MAESTRO_API_KEY    string = "bODB5wkcG0EGkBgFRHPkSgqQQ186WoTB"
	// MAESTRO_API_KEY           string = "PsK0OfakoWpsbIbtAsV9edNVZ2Edy3aZ"
	KUPO_ENDPOINT             string = "http://localhost:1442"
	BLINKLABS_KUPO_ENDPOINT   string = "https://preprod-kupo.blinklabs.io"
	BLINKLABS_OGMIOS_ENDPOINT string = "wss://preprod-ogmios.blinklabs.io"

	EKIVAL_FEE int = 25_000_000
)

const (
	ORDER_ID    string = "AwesomeID_TT_RR_1"
	TO_LOVELACE int64  = 1000000
	// ORDER_AMOUNT int64  = 854678691
	ORDER_AMOUNT      int64 = 501 * TO_LOVELACE
	PRECISION         int64 = 10
	MIN_COLLATERAL    int64 = 25 * TO_LOVELACE
	MAKER_MIN_FEE     int64 = 1250000
	TAKER_MIN_FEE     int64 = 3750000
	CANCEL_MIN_FEE    int64 = 3 * TO_LOVELACE
	MIN_ORDER_AMOUNT  int64 = 10 * TO_LOVELACE
	ORDER_THRESHOLD   int64 = 500 * TO_LOVELACE
	CANCEL_PENALTY    int64 = 2 * CANCEL_MIN_FEE
	ADA_COLLATERAL    int64 = 3 * TO_LOVELACE
	EXTENDING_PENALTY int64 = 5 * TO_LOVELACE
	DISPUTE_MIN_FEE   int64 = 5 * TO_LOVELACE
	DISPUTE_PENALTY   int64 = 2 * DISPUTE_MIN_FEE

	MAKER_DEADLINE int64 = 7200
	TAKER_DEADLINE int64 = 3600
)

const (
	TAKER_COMMIT_TO_ORDER_UTXO_TXID       string = "82e6ceede062175b547ee92c3f6faea5505cdc2d7dec664a0d927468955d8252"
	TAKER_COMMIT_TO_ORDER_UTXO_TXID_INDEX int    = 0

	TAKER_CONFIRM_REMIT_UTXO_TXID       = "aa45a5579783c19d2a14621cddbcf9af1f0e64dbb2e9bc78e628070dd6ce9b9d"
	TAKER_CONFIRM_REMIT_UTXO_TXID_INDEX = 0

	MAKER_CONFIRM_RECEIPT_UTXO_TXID       string = "f2753cbd3147be4a1d523896d70e30885eeba5bc7d18f8cc829a28759584d866"
	MAKER_CONFIRM_RECEIPT_UTXO_TXID_INDEX int    = 0

	THE_MAKER_DEADLINE int64 = 1763768995761
	THE_TAKER_DEADLINE int64 = 1763765395761
)
