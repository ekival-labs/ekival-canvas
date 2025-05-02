package viewmodel

import (
	"fmt"

	"github.com/ekival-labs/ekival-canvas/model"
)

type UserTxInfo struct {
	Address            string        `json:"address"`
	ChangeAddress      string        `json:"changeAddress"`
	UserUTxOs          []model.EUTxO `json:"UserUTxOs"`
	CollateralUTxO     model.EUTxO   `json:"CollateralUTxO,omitempty"`
	OrderUtxo          model.EUTxO   `json:"OrderUtxo,omitempty"`
	TradeTokenName     string        `json:"TradeTokenName,omitempty"`
	TradeTokenPolicyId string        `json:"TradeTokenPolicyId,omitempty"`
	OrderId            string        `json:"OrderId,omitempty"`
	OrderAmount        int64         `json:"OrderAmount,omitempty"`
	MakerAddress       string        `json:"MakerAddress,omitempty"`
	MakerRepAddress    string        `json:"MakerRepAddress,omitempty"`
	TakerAddress       string        `json:"TakerAddress,omitempty"`
	TakerRepAddress    string        `json:"TakerRepAddress,omitempty"`
	MakerDeadline      int64         `json:"MakerDeadline,omitempty"`
	TakerDeadline      int64         `json:"TakerDeadline,omitempty"`
	Precision          int64         `json:"Precision,omitempty"`
	CollateralPct      int64         `json:"CollateralPct,omitempty"`
	MakerPct           int64         `json:"MakerPct,omitempty"`
	TakerPct           int64         `json:"TakerPct,omitempty"`
	CancelPct          int64         `json:"CancelPct,omitempty"`
	MinCollateral      int64         `json:"MinCollateral,omitempty"`
	MakerMinFee        int64         `json:"MakerMinFee,omitempty"`
	TakerMinFee        int64         `json:"TakerMinFee,omitempty"`
	CancelMinFee       int64         `json:"CancelMinFee,omitempty"`
	MinOrderAmount     int64         `json:"MinOrderAmount,omitempty"`
	OrderThreshold     int64         `json:"OrderThreshold,omitempty"`
	CancelPenalty      int64         `json:"CancelPenalty,omitempty"`
	AdaCollateral      int64         `json:"AdaCollateral,omitempty"`
}

type TxResponse struct {
	TxCBOR string `json:"txCBOR"`
	TxID   string `json:"txID"`
}

func (u *UserTxInfo) IsValid() error {
	if u.Address == "" || u.ChangeAddress == "" || u.UserUTxOs == nil {
		return fmt.Errorf("error at IsValid: one or more fields are empty")
	}
	return nil
}

func (t *TxResponse) IsValid() error {
	if t.TxCBOR == "" {
		return fmt.Errorf("error at IsValid: one or more fields are empty")
	}
	return nil
}
