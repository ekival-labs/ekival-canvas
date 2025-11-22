package utility

import (
	"errors"
)

func CalculateFee(precision, orderThreshold, minOrder, minFee, feePercentage, orderAmount int64) (int64, error) {
	if orderAmount <= orderThreshold && orderAmount >= minOrder {
		return minFee, nil
	} else if orderAmount > orderThreshold {
		return (orderAmount * precision) / feePercentage, nil
	} else {
		return 0, errors.New("Order amount does not meet criteria")
	}
}
