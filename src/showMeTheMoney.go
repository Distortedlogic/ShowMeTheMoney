// package with only one function GetVestingBalance

package showMeTheMoney

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// Returns the vesting coins as sdk.Coins for a periodic Vesting Acct
func GetVestingBalance(
	account vestingtypes.PeriodicVestingAccount,
	currentUnixTime int64,
) sdk.Coins {
	vesting := account.OriginalVesting // Start at original vesting
	if currentUnixTime < account.StartTime {
		return vesting // before start time, just return original vesting
	}
	ct := account.StartTime
	periods := account.VestingPeriods
	for _, period := range periods {
		if currentUnixTime-ct < period.Length {
			break // next period wont be vested so break loop
		}
		vesting = vesting.Sub(period.Amount) // substract the vested coins from the vesting coins
		ct += period.Length                  // accum period lengths onto start time
	}
	return vesting // ShowThemTheMoney!
}
