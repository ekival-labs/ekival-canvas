package viewmodel

import (
	"fmt"
)

type EUTxO struct {
	TxID      string `json:"TxID"`
	TxIDIndex int    `json:"TxIDIndex"`
}

type UserTxInfo struct {
	Address        string  `json:"address"`
	ChangeAddress  string  `json:"changeAddress"`
	UserUTxOs      []EUTxO `json:"UserUTxOs"`
	CollateralUTxO EUTxO   `json:"CollateralUTxO,omitempty"`
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
	if t.TxCBOR == "" || t.TxID == "" {
		return fmt.Errorf("error at IsValid: one or more fields are empty")
	}
	return nil
}
