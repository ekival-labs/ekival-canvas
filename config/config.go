package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	PayloadExpireAtMinute time.Duration = 15
)

var (
	GlobalConfig *Config
)

type Config struct {
	ContractsController SpendingContractConfig `json:"ContractsController"`
	MainTreasury        SpendingContractConfig `json:"MainTreasury"`
	OfferAdmin          string                 `json:"OfferAdmin"`
	ADAMarketplace      ADAMarketplace         `json:"ada"`
	TMoneyMarketplace   TokenMarketplace       `json:"tMoney"`
}

type ADAMarketplace struct {
	AdaP2PBuyEscrow       SpendingContractConfig `json:"AdaP2PBuyEscrow"`
	APBST                 MintingContractConfig  `json:"APBST"`
	AdaP2PSellEscrow      SpendingContractConfig `json:"AdaP2PSellEscrow"`
	APSST                 MintingContractConfig  `json:"APSST"`
	AdminWalletsMnemonics AdminsMnemonics        `json:"AdminWalletsMnemonics"`
}
type TokenMarketplace struct {
	TokenP2PBuyEscrow     SpendingContractConfig `json:"TokenP2PBuyEscrow"`
	TPBST                 MintingContractConfig  `json:"TPBST"`
	TokenP2PSellEscrow    SpendingContractConfig `json:"TokenP2PSellEscrow"`
	TPSST                 MintingContractConfig  `json:"TPSST"`
	AdminWalletsMnemonics AdminsMnemonics        `json:"AdminWalletsMnemonics"`
}

type AdminsMnemonics struct {
	BuyAdmin  string `json:"BuyAdmin"`
	SellAdmin string `json:"SellAdmin"`
	BSTAdmin  string `json:"BSTAdmin"`
	SSTAdmin  string `json:"SSTAdmin"`
}
type SpendingContractConfig struct {
	Address  string `json:"Address"`
	RefTxID  string `json:"RefTxID"`
	RefTxIDx int    `json:"RefTxIDx"`
}

type MintingContractConfig struct {
	PolicyID string `json:"PolicyID"`
	RefTxID  string `json:"RefTxID"`
	RefTxIDx int    `json:"RefTxIDx"`
}

func Load(configFile string) error {

	if configFile != "" {
		buf, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("error reading config file: %s", err)
		}
		err = json.Unmarshal(buf, &GlobalConfig)
		if err != nil {
			return fmt.Errorf("error parsing config file: %v", err)
		}
	}

	return nil
}

func GetGlobalConfig() *Config {
	return GlobalConfig
}
