package config

import (
	"encoding/hex"
	"strings"

	"ekival-canvas/constants"

	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/serialization/Key"
	"github.com/blinklabs-io/bursa"
)

type Wallet struct {
	Address Address.Address
	PKH     serialization.PubKeyHash
	Vkey    Key.VerificationKey
	Skey    Key.SigningKey
}

var (
	offerAdminWallet         = &Wallet{}
	adaP2PBuyAdminWallet     = &Wallet{}
	adaP2PSellAdminWallet    = &Wallet{}
	tMoneyP2PBuyAdminWallet  = &Wallet{}
	tMoneyP2PSellAdminWallet = &Wallet{}
	aPBSTAdminWallet         = &Wallet{}
	aPSSTAdminWallet         = &Wallet{}
	tMoneyBSTAdminWallet     = &Wallet{}
	tMoneySSTAdminWallet     = &Wallet{}
)

func WalletSetup() {

	cfg := GetGlobalConfig()

	offerAdminWallet = SetWallet(toMnemonic(cfg.OfferAdmin))
	adaP2PBuyAdminWallet = SetWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.BuyAdmin))
	adaP2PSellAdminWallet = SetWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.SellAdmin))
	aPBSTAdminWallet = SetWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.BSTAdmin))
	aPSSTAdminWallet = SetWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.SSTAdmin))
	tMoneyP2PBuyAdminWallet = SetWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.BuyAdmin))
	tMoneyP2PSellAdminWallet = SetWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.SellAdmin))
	tMoneyBSTAdminWallet = SetWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.BSTAdmin))
	tMoneySSTAdminWallet = SetWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.SSTAdmin))
}

func SetWallet(mnemonic string) *Wallet {
	rootKey, err := bursa.GetRootKeyFromMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}
	accountKey := bursa.GetAccountKey(rootKey, 0)
	paymentKey := bursa.GetPaymentKey(accountKey, 0)
	address, err := Address.DecodeAddress(bursa.GetAddress(accountKey, constants.NETWORK, 0).String())
	if err != nil {
		panic(err)
	}
	vKeyBytes, err := hex.DecodeString(bursa.GetPaymentVKey(paymentKey).CborHex)
	if err != nil {
		panic(err)

	}
	sKeyBytes, err := hex.DecodeString(bursa.GetPaymentSKey(paymentKey).CborHex)
	if err != nil {
		panic(err)
	}
	vKeyBytes = vKeyBytes[2:]
	sKeyBytes = sKeyBytes[2:]
	sKeyBytes = append(sKeyBytes[:64], sKeyBytes[96:]...)

	return &Wallet{
		Address: address,
		PKH:     serialization.PubKeyHash(address.PaymentPart),
		Vkey:    Key.VerificationKey{Payload: vKeyBytes},
		Skey:    Key.SigningKey{Payload: sKeyBytes},
	}
}

func toMnemonic(seedPhrase string) (mnemonic string) {
	words := strings.Fields(seedPhrase)
	mnemonic = strings.Join(words, " ")
	return mnemonic
}

func GetOfferAdminWallet() *Wallet {
	return offerAdminWallet
}

func GetAdaP2PBuyAdminWallet() *Wallet {
	return adaP2PBuyAdminWallet
}

func GetAdaP2PSellAdminWallet() *Wallet {
	return adaP2PSellAdminWallet
}

func GetTMoneyP2PBuyAdminWallet() *Wallet {
	return tMoneyP2PBuyAdminWallet
}

func GetTMoneyP2PSellAdminWallet() *Wallet {
	return tMoneyP2PSellAdminWallet
}

func GetAPBSTAdminWallet() *Wallet {
	return aPBSTAdminWallet
}

func GetAPSSTAdminWallet() *Wallet {
	return aPSSTAdminWallet
}

func GetTMoneyBSTAdminWallet() *Wallet {
	return tMoneyBSTAdminWallet
}

func GetTMoneySSTAdminWallet() *Wallet {
	return tMoneySSTAdminWallet
}
