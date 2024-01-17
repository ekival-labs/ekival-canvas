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
	AdaP2PBuyEscrow       SpendingContractConfig `json:"AdaP2PBuyEscrow"`
	APBST                 MintingContractConfig  `json:"APBST"`
	AdaP2PSellEscrow      SpendingContractConfig `json:"AdaP2PSellEscrow"`
	APSST                 MintingContractConfig  `json:"APSST"`
	EkiP2PBuyEscrow       SpendingContractConfig `json:"EkiP2PBuyEscrow"`
	EKI_BST               MintingContractConfig  `json:"EKI_BST"`
	EkiP2PSellEscrow      SpendingContractConfig `json:"EkiP2PSellEscrow"`
	EKI_SST               MintingContractConfig  `json:"EKI_SST"`
	Staking               SpendingContractConfig `json:"Staking"`
	MainTreasury          SpendingContractConfig `json:"MainTreasury"`
	AdminWalletsMnemonics AdminsMnemonics        `json:"AdminWalletsMnemonics"`
}

type AdminsMnemonics struct {
	OfferAdmin      string `json:"OfferAdmin"`
	AdaP2PBuyAdmin  string `json:"AdaP2PBuyAdmin"`
	AdaP2PSellAdmin string `json:"AdaP2PSellAdmin"`
	EkiP2PBuyAdmin  string `json:"EkiP2PBuyAdmin"`
	EkiP2PSellAdmin string `json:"EkiP2PSellAdmin"`
	APBSTAdmin      string `json:"APBSTAdmin"`
	APSSTAdmin      string `json:"APSSTAdmin"`
	EKI_BSTAdmin    string `json:"EKI_BSTAdmin"`
	EKI_SSTAdmin    string `json:"EKI_SSTAdmin"`
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
