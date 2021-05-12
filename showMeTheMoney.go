package showMeTheMoney

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

func GetVestingBalance(
	account vestingtypes.PeriodicVestingAccount,
	currentUnixTime int64,
) sdk.Coins {
	vesting := account.OriginalVesting
	if currentUnixTime < account.StartTime {
		return vesting
	}
	ct := account.StartTime
	periods := account.VestingPeriods
	for _, period := range periods {
		if currentUnixTime-ct < period.Length {
			break
		}
		vesting = vesting.Sub(period.Amount)
		ct += period.Length
	}
	return vesting
}
