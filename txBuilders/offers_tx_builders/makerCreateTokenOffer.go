package offers_tx_builders

import (
	"ekival-canvas/config"
	"ekival-canvas/model"
	"ekival-canvas/utility"
	"encoding/hex"
	"log"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/txBuilding/Utils"
)

func MakerCreateTokenOffer(offer *model.TokenOfferTxInfo, treasury *model.TreasuryInfo, adminWallet *config.AdminWallet) (string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			return
		}
	}()

	apolloBE := apollo.New(&config.BFC)
	// apolloBE = apolloBE.SetWalletFromAddress(offer.MakerAddress)
	apolloBE = apolloBE.SetWalletFromBech32(offer.MakerAddress.String())

	userUtxos, err := utility.GetUserUTxOs(offer.UserUtxos)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	collateralUtxo := config.BFC.GetUtxoFromRef(offer.CollateralUtxo.TxID, offer.CollateralUtxo.TxIDIndex)

	lastSlot := config.BFC.LastBlockSlot()

	apolloBE, err = apolloBE.
		AddLoadedUTxOs(userUtxos...).
		SetChangeAddress(offer.ChangeAddress).
		AddCollateral(*collateralUtxo).
		PayToContract(
			treasury.Address,
			treasury.Datum,
			offer.EkivalFeeLovelace,
			true,
			apollo.Unit{
				PolicyId: offer.TradeTokenPolicyId,
				Name:     offer.TradeTokenName,
				Quantity: offer.EkivalFeeToken,
			},
		).
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(offer.MakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		Complete()

	if err != nil {
		log.Println(err)
		return "", "", err
	}

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	tx := apolloBE.GetTx()

	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return Utils.ToCbor(tx), hex.EncodeToString(txHash), nil

}
