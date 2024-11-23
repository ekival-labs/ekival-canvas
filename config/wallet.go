package config

import (
	"ekival-canvas/constants"
	"encoding/hex"
	"strings"

	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/serialization/Key"
	"github.com/blinklabs-io/bursa"
)

type AdminWallet struct {
	AdminPKH  serialization.PubKeyHash
	AdminVkey Key.VerificationKey
	AdminSkey Key.SigningKey
}

var (
	offerAdminWallet         = &AdminWallet{}
	adaP2PBuyAdminWallet     = &AdminWallet{}
	adaP2PSellAdminWallet    = &AdminWallet{}
	tMoneyP2PBuyAdminWallet  = &AdminWallet{}
	tMoneyP2PSellAdminWallet = &AdminWallet{}
	aPBSTAdminWallet         = &AdminWallet{}
	aPSSTAdminWallet         = &AdminWallet{}
	tMoneyBSTAdminWallet     = &AdminWallet{}
	tMoneySSTAdminWallet     = &AdminWallet{}
)

func WalletSetup() {

	cfg := GetGlobalConfig()

	offerAdminWallet = SetAdminWallet(toMnemonic(cfg.OfferAdmin))
	adaP2PBuyAdminWallet = SetAdminWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.BuyAdmin))
	adaP2PSellAdminWallet = SetAdminWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.SellAdmin))
	aPBSTAdminWallet = SetAdminWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.BSTAdmin))
	aPSSTAdminWallet = SetAdminWallet(toMnemonic(cfg.ADAMarketplace.AdminWalletsMnemonics.SSTAdmin))
	tMoneyP2PBuyAdminWallet = SetAdminWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.BuyAdmin))
	tMoneyP2PSellAdminWallet = SetAdminWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.SellAdmin))
	tMoneyBSTAdminWallet = SetAdminWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.BSTAdmin))
	tMoneySSTAdminWallet = SetAdminWallet(toMnemonic(cfg.TMoneyMarketplace.AdminWalletsMnemonics.SSTAdmin))
}

func SetAdminWallet(mnemonic string) *AdminWallet {
	rootKey, err := bursa.GetRootKeyFromMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}
	accountKey := bursa.GetAccountKey(rootKey, 0)
	paymentKey := bursa.GetPaymentKey(accountKey, 0)
	adminAddress, err := Address.DecodeAddress(bursa.GetAddress(accountKey, constants.NETWORK, 0).String())
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

	return &AdminWallet{
		AdminPKH:  serialization.PubKeyHash(adminAddress.PaymentPart),
		AdminVkey: Key.VerificationKey{Payload: vKeyBytes},
		AdminSkey: Key.SigningKey{Payload: sKeyBytes},
	}
}

func toMnemonic(seedPhrase string) (mnemonic string) {
	words := strings.Fields(seedPhrase)
	mnemonic = strings.Join(words, " ")
	return mnemonic
}

func GetOfferAdminWallet() *AdminWallet {
	return offerAdminWallet
}

func GetAdaP2PBuyAdminWallet() *AdminWallet {
	return adaP2PBuyAdminWallet
}

func GetAdaP2PSellAdminWallet() *AdminWallet {
	return adaP2PSellAdminWallet
}

func GetTMoneyP2PBuyAdminWallet() *AdminWallet {
	return tMoneyP2PBuyAdminWallet
}

func GetTMoneyP2PSellAdminWallet() *AdminWallet {
	return tMoneyP2PSellAdminWallet
}

func GetAPBSTAdminWallet() *AdminWallet {
	return aPBSTAdminWallet
}

func GetAPSSTAdminWallet() *AdminWallet {
	return aPSSTAdminWallet
}

func GetTMoneyBSTAdminWallet() *AdminWallet {
	return tMoneyBSTAdminWallet
}

func GetTMoneySSTAdminWallet() *AdminWallet {
	return tMoneySSTAdminWallet
}
