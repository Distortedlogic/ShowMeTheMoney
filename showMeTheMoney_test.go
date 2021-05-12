package showMeTheMoney

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// Creating Base Account for testing from Bech32 address and pubKey
var address, _ = sdk.AccAddressFromBech32(
	"cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
)
var pubKey, _ = sdk.GetPubKeyFromBech32(
	sdk.Bech32PubKeyTypeAccPub,
	"cosmospub1addwnpepqd5xvvdrw7dsfe89pcr9amlnvx9qdkjgznkm2rlfzesttpjp50jy2lueqp2",
)
var bassAcc = authtypes.NewBaseAccount(address, pubKey, 13748, 240)

func TestGetVestCoins(t *testing.T) {
	testCases := []struct {
		name             string
		currentTimestamp int64
		startTime        int64
		vestingPeriods   vestingtypes.Periods
		originalVesting  sdk.Coins
		expected         sdk.Coins
	}{
		{
			"Current time less than start time",
			time.Now().Unix(),
			time.Now().AddDate(0, 0, 1).Unix(),
			vestingtypes.Periods{
				{
					Length: 6763707,
					Amount: sdk.Coins{{Denom: "ukava", Amount: sdk.NewInt(21040)}},
				},
			},
			sdk.Coins{
				{Denom: "hard", Amount: sdk.NewInt(47513838)},
				{Denom: "ukava", Amount: sdk.NewInt(10749247)},
			},
			sdk.Coins{
				{Denom: "hard", Amount: sdk.NewInt(47513838)},
				{Denom: "ukava", Amount: sdk.NewInt(10749247)},
			},
		},
		{
			"Should be vested",
			time.Now().Unix(),
			time.Now().AddDate(0, 0, -1).Unix(),
			vestingtypes.Periods{
				{
					Length: 60 * 60 * 12,
					Amount: sdk.Coins{{Denom: "ukava", Amount: sdk.NewInt(21040)}},
				},
			},
			sdk.Coins{
				{Denom: "hard", Amount: sdk.NewInt(47513838)},
				{Denom: "ukava", Amount: sdk.NewInt(10749247)},
			},
			sdk.Coins{
				{Denom: "hard", Amount: sdk.NewInt(47513838)},
				{Denom: "ukava", Amount: sdk.NewInt(10749247 - 21040)},
			},
		},
		{
			"Should not be vested",
			time.Now().Unix(),
			time.Now().AddDate(0, 0, -1).Unix(),
			vestingtypes.Periods{
				{
					Length: 60 * 60 * 24 * 2,
					Amount: sdk.Coins{{Denom: "ukava", Amount: sdk.NewInt(21040)}},
				},
			},
			sdk.Coins{
				{Denom: "hard", Amount: sdk.NewInt(47513838)},
				{Denom: "ukava", Amount: sdk.NewInt(10749247)},
			},
			sdk.Coins{
				{Denom: "hard", Amount: sdk.NewInt(47513838)},
				{Denom: "ukava", Amount: sdk.NewInt(10749247)},
			},
		},
	}
	for idx, tc := range testCases {
		newPVA := vestingtypes.NewPeriodicVestingAccount(
			bassAcc,
			tc.originalVesting,
			tc.startTime,
			tc.vestingPeriods,
		)
		if vested := GetVestingBalance(*newPVA, tc.currentTimestamp); !vested.IsEqual(tc.expected) {
			t.Errorf("Test %d: %s failed - output %s expected %s", idx, tc.name, vested, tc.expected)
		}
	}

}
