package tests

import (
	"ekival-canvas/utility"
	"testing"
)

func TestCalculateFee(t *testing.T) {
	var (
		toLovelace int64 = 1000000

		Precision      int64 = 10
		MakerPct       int64 = utility.ToFraction(0.25)
		TakerPct       int64 = utility.ToFraction(0.75)
		MakerMinFee    int64 = 1250000
		TakerMinFee    int64 = 3750000
		OrderThreshold int64 = 500 * toLovelace
		MinOrder       int64 = 10 * toLovelace
	)

	testCases := []struct {
		name           string
		orderAmount    int64
		orderThreshold int64
		feePercentage  int64
		minFee         int64
		expectedFee    int64
		minOrder       int64
		err            error
	}{
		{
			name:           "Maker Fee Below Threshold",
			orderAmount:    100 * toLovelace,
			orderThreshold: OrderThreshold,
			feePercentage:  MakerPct,
			minFee:         MakerMinFee,
			expectedFee:    MakerMinFee,
			minOrder:       MinOrder,
			err:            nil,
		},
		{
			name:           "Maker Fee Above Threshold",
			orderAmount:    1000 * toLovelace,
			orderThreshold: OrderThreshold,
			feePercentage:  MakerPct,
			minFee:         MakerMinFee,
			expectedFee:    ((1000 * toLovelace * Precision) / MakerPct),
			minOrder:       MinOrder,
			err:            nil,
		},
		{
			name:           "Taker Fee Below Threshold",
			orderAmount:    100 * toLovelace,
			orderThreshold: OrderThreshold,
			feePercentage:  TakerPct,
			minFee:         TakerMinFee,
			expectedFee:    TakerMinFee,
			minOrder:       MinOrder,
			err:            nil,
		},
		{
			name:           "Taker Fee Above Threshold",
			orderAmount:    1000 * toLovelace,
			orderThreshold: OrderThreshold,
			feePercentage:  TakerPct,
			minFee:         TakerMinFee,
			expectedFee:    ((1000 * toLovelace * Precision) / TakerPct),
			minOrder:       MinOrder,
			err:            nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running test case: %s", tc.name)
			t.Logf("  Order Amount: %d, Order Threshold: %d, Fee Percentage: %d, Min Fee: %d", tc.orderAmount, tc.orderThreshold, tc.feePercentage, tc.minFee)
			actualFee, err := utility.CalculateFee(Precision, tc.orderThreshold, tc.minOrder, tc.minFee, tc.feePercentage, tc.orderAmount)
			if err != nil {
				t.Errorf("For %s: Error: %v", tc.name, err)
			}
			t.Logf("  Expected Fee: %d, Actual Fee: %d", tc.expectedFee, actualFee)
			if actualFee != tc.expectedFee {
				t.Errorf("For %s: Expected fee %d, but got %d", tc.name, tc.expectedFee, actualFee)
			}
		})
	}
}
