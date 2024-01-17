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
	offerAdminWallet      = &AdminWallet{}
	adaP2PBuyAdminWallet  = &AdminWallet{}
	adaP2PSellAdminWallet = &AdminWallet{}
	EkiP2PBuyAdminWallet  = &AdminWallet{}
	EkiP2PSellAdminWallet = &AdminWallet{}
	aPBSTAdminWallet      = &AdminWallet{}
	aPSSTAdminWallet      = &AdminWallet{}
	EKI_BSTAdminWallet    = &AdminWallet{}
	EKI_SSTAdminWallet    = &AdminWallet{}
)

func WalletSetup() {

	cfg := GetGlobalConfig()

	offerAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.OfferAdmin))
	adaP2PBuyAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.AdaP2PBuyAdmin))
	adaP2PSellAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.AdaP2PSellAdmin))
	EkiP2PBuyAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.EkiP2PBuyAdmin))
	EkiP2PSellAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.EkiP2PSellAdmin))
	aPBSTAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.APBSTAdmin))
	aPSSTAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.APSSTAdmin))
	EKI_BSTAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.EKI_BSTAdmin))
	EKI_SSTAdminWallet = SetAdminWallet(toMnemonic(cfg.AdminWalletsMnemonics.EKI_SSTAdmin))
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

func GetEkiP2PBuyAdminWallet() *AdminWallet {
	return EkiP2PBuyAdminWallet
}

func GetEkiP2PSellAdminWallet() *AdminWallet {
	return EkiP2PSellAdminWallet
}

func GetAPBSTAdminWallet() *AdminWallet {
	return aPBSTAdminWallet
}

func GetAPSSTAdminWallet() *AdminWallet {
	return aPSSTAdminWallet
}

func GetEKI_BSTAdminWallet() *AdminWallet {
	return EKI_BSTAdminWallet
}

func GetEKI_SSTAdminWallet() *AdminWallet {
	return EKI_SSTAdminWallet
}
