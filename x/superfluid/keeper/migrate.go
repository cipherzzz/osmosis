package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cl "github.com/osmosis-labs/osmosis/v14/x/concentrated-liquidity"
	cltypes "github.com/osmosis-labs/osmosis/v14/x/concentrated-liquidity/types"
	gammtypes "github.com/osmosis-labs/osmosis/v14/x/gamm/types"
)

// UnlockAndMigrate unlocks a balancer pool lock, exits the pool and migrates the LP position to a full range concentrated liquidity position.
// If the lock is also superfluid delegated, it will undelegate the superfluid position.
func (k Keeper) UnlockAndMigrate(ctx sdk.Context, sender sdk.AccAddress, lockId uint64, sharesToMigrate sdk.Coin) (amount0, amount1 sdk.Int, liquidity sdk.Dec, poolIdLeaving, poolIdEntering uint64, frozenUntil time.Time, err error) {
	// Get the balancer poolId by parsing the gamm share denom.
	poolIdLeaving = gammtypes.MustGetPoolIdFromShareDenom(sharesToMigrate.Denom)

	// Ensure a governance sanctioned link exists between the balancer pool and the concentrated pool.
	poolIdEntering, err = k.gk.GetLinkedConcentratedPool(ctx, poolIdLeaving)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, err
	}

	// Get the concentrated pool from the provided ID and type cast it to ConcentratedPoolExtension.
	poolI, err := k.clk.GetPool(ctx, poolIdEntering)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, err
	}
	concentratedPool, ok := poolI.(cltypes.ConcentratedPoolExtension)
	if !ok {
		// If the conversion fails, return an error.
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, fmt.Errorf("given pool does not implement ConcentratedPoolExtension, implements %T", poolI)
	}

	// Consistency check that lockID corresponds to sender, and contains correct LP shares.
	// These are expected to be true by the caller, but good to double check
	lock, err := k.validateLockForUnpool(ctx, sender, poolIdLeaving, lockId)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, err
	}

	// Before we break the lock, we must note the time remaining on the lock.
	// We will be freezing the concentrated liquidity position for this duration.
	freezeDuration := k.getExistingLockRemainingDuration(ctx, lock)

	// If superfluid delegated, superfluid undelegate
	// This also burns the underlying synthetic osmo
	err = k.unbondSuperfluidIfExists(ctx, sender, lockId)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, err
	}

	// Finish unlocking directly for locked locks
	// this also unlocks locks that were in the unlocking queue
	err = k.lk.ForceUnlock(ctx, *lock)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, err
	}

	// Exit the balancer pool position.
	exitCoins, err := k.gk.ExitPool(ctx, sender, poolIdLeaving, sharesToMigrate.Amount, sdk.NewCoins())
	if err != nil {
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, err
	}

	// Determine the max and min ticks for the concentrated pool we are migrating to.
	minTick, maxTick := cl.GetMinAndMaxTicksFromExponentAtPriceOne(concentratedPool.GetPrecisionFactorAtPriceOne())

	frozenUntil = ctx.BlockTime().Add(freezeDuration)

	// Create a full range (min to max tick) concentrated liquidity position.
	amount0, amount1, liquidity, err = k.clk.CreatePosition(ctx, poolIdEntering, sender, exitCoins.AmountOf(concentratedPool.GetToken0()), exitCoins.AmountOf(concentratedPool.GetToken1()), sdk.ZeroInt(), sdk.ZeroInt(), minTick, maxTick, frozenUntil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, sdk.Dec{}, 0, 0, time.Time{}, err
	}

	return amount0, amount1, liquidity, poolIdLeaving, poolIdEntering, frozenUntil, nil
}
