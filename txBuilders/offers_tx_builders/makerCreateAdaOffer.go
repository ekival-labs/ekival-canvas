package offers_tx_builders

import (
	"ekival-canvas/config"
	"ekival-canvas/model"
	"ekival-canvas/utility"
	"encoding/hex"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/txBuilding/Utils"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

func MakerCreateAdaOffer(offer *model.AdaOfferTxInfo, treasury *model.TreasuryInfo, adminWallet *config.AdminWallet) (string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			fiberLogger.Panic("Panic occurred: %v", err)
			return
		}
	}()

	apolloBE := apollo.New(&config.CHAIN_CTX)
	// apolloBE = apolloBE.SetWalletFromAddress(offer.MakerAddress)
	apolloBE = apolloBE.SetWalletFromBech32(offer.MakerAddress.String())

	userUtxos, err := utility.GetUserUTxOs(offer.UserUtxos)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	collateralUtxo, err := config.CHAIN_CTX.GetUtxoFromRef(offer.CollateralUtxo.TxID, offer.CollateralUtxo.TxIDIndex)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	lastSlot, err := config.CHAIN_CTX.LastBlockSlot()
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	apolloBE, err = apolloBE.
		AddLoadedUTxOs(userUtxos...).
		AddCollateral(*collateralUtxo).
		SetChangeAddress(offer.ChangeAddress).
		PayToContract(
			treasury.Address, treasury.Datum, offer.EkivalFeeLovelace, true,
		).
		AddRequiredSigner(adminWallet.AdminPKH).
		AddRequiredSigner(serialization.PubKeyHash(offer.MakerAddress.PaymentPart)).
		SetTtl(int64(lastSlot) + 300).
		Complete()

	fiberLogger.Error(serialization.PubKeyHash(adminWallet.AdminPKH))

	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	apolloBE, err = apolloBE.SignWithSkey(adminWallet.AdminVkey, adminWallet.AdminSkey)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	tx := apolloBE.GetTx()
	txHash, err := tx.TransactionBody.Hash()
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	cbor, err := Utils.ToCbor(tx)
	if err != nil {
		fiberLogger.Error(err)
		return "", "", err
	}

	return cbor, hex.EncodeToString(txHash), nil

}
