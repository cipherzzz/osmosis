package keeper

import (
	"fmt"
	"strconv"

	"github.com/osmosis-labs/osmosis/v14/x/protorev/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/osmosis-labs/osmosis/osmoutils"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------- Trading Stores  ---------------------- //

// GetTokenPairArbRoutes returns the token pair arb routes given two denoms
func (k Keeper) GetTokenPairArbRoutes(ctx sdk.Context, tokenA, tokenB string) (*types.TokenPairArbRoutes, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairRoutes)
	key := types.GetKeyPrefixRouteForTokenPair(tokenA, tokenB)

	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, fmt.Errorf("no routes found for token pair %s-%s", tokenA, tokenB)
	}

	tokenPairArbRoutes := &types.TokenPairArbRoutes{}
	err := tokenPairArbRoutes.Unmarshal(bz)
	if err != nil {
		return nil, err
	}

	return tokenPairArbRoutes, nil
}

// GetAllTokenPairArbRoutes returns all the token pair arb routes
func (k Keeper) GetAllTokenPairArbRoutes(ctx sdk.Context) ([]*types.TokenPairArbRoutes, error) {
	routes := make([]*types.TokenPairArbRoutes, 0)

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixTokenPairRoutes)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		tokenPairArbRoutes := &types.TokenPairArbRoutes{}
		err := tokenPairArbRoutes.Unmarshal(iterator.Value())
		if err != nil {
			return nil, err
		}

		routes = append(routes, tokenPairArbRoutes)
	}

	return routes, nil
}

// SetTokenPairArbRoutes sets the token pair arb routes given two denoms
func (k Keeper) SetTokenPairArbRoutes(ctx sdk.Context, tokenA, tokenB string, tokenPair *types.TokenPairArbRoutes) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairRoutes)
	key := types.GetKeyPrefixRouteForTokenPair(tokenA, tokenB)

	bz, err := tokenPair.Marshal()
	if err != nil {
		return err
	}

	store.Set(key, bz)

	return nil
}

// DeleteAllTokenPairArbRoutes deletes all the token pair arb routes
func (k Keeper) DeleteAllTokenPairArbRoutes(ctx sdk.Context) {
	k.DeleteAllEntriesForKeyPrefix(ctx, types.KeyPrefixTokenPairRoutes)
}

// GetAllBaseDenoms returns all of the base denoms (sorted by priority in descending order) used to build cyclic arbitrage routes
func (k Keeper) GetAllBaseDenoms(ctx sdk.Context) ([]*types.BaseDenom, error) {
	baseDenoms := make([]*types.BaseDenom, 0)

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixBaseDenoms)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		baseDenom := &types.BaseDenom{}
		err := baseDenom.Unmarshal(iterator.Value())
		if err != nil {
			return []*types.BaseDenom{}, err
		}

		baseDenoms = append(baseDenoms, baseDenom)
	}

	return baseDenoms, nil
}

// SetBaseDenoms sets all of the base denoms used to build cyclic arbitrage routes. The base denoms priority
// order is going to match the order of the base denoms in the slice.
func (k Keeper) SetBaseDenoms(ctx sdk.Context, baseDenoms []*types.BaseDenom) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBaseDenoms)

	for i, baseDenom := range baseDenoms {
		key := types.GetKeyPrefixBaseDenom(uint64(i))

		bz, err := baseDenom.Marshal()
		if err != nil {
			return err
		}
		store.Set(key, bz)
	}

	return nil
}

// DeleteBaseDenoms deletes all of the base denoms
func (k Keeper) DeleteBaseDenoms(ctx sdk.Context) {
	k.DeleteAllEntriesForKeyPrefix(ctx, types.KeyPrefixBaseDenoms)
}

// GetPoolForDenomPair returns the id of the highest liquidty pool between the base denom and the denom to match
func (k Keeper) GetPoolForDenomPair(ctx sdk.Context, baseDenom, denomToMatch string) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDenomPairToPool)
	key := types.GetKeyPrefixDenomPairToPool(baseDenom, denomToMatch)

	bz := store.Get(key)
	if len(bz) == 0 {
		return 0, fmt.Errorf("highest liquidity pool between base %s and match denom %s not found", baseDenom, denomToMatch)
	}

	poolId := sdk.BigEndianToUint64(bz)
	return poolId, nil
}

// SetPoolForDenomPair sets the id of the highest liquidty pool between the base denom and the denom to match
func (k Keeper) SetPoolForDenomPair(ctx sdk.Context, baseDenom, denomToMatch string, poolId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDenomPairToPool)
	key := types.GetKeyPrefixDenomPairToPool(baseDenom, denomToMatch)

	store.Set(key, sdk.Uint64ToBigEndian(poolId))
}

// DeleteAllPoolsForBaseDenom deletes all the pools for the given base denom
func (k Keeper) DeleteAllPoolsForBaseDenom(ctx sdk.Context, baseDenom string) {
	key := append(types.KeyPrefixDenomPairToPool, types.GetKeyPrefixDenomPairToPool(baseDenom, "")...)
	k.DeleteAllEntriesForKeyPrefix(ctx, key)
}

// DeleteAllEntriesForKeyPrefix deletes all the entries from the store for the given key prefix
func (k Keeper) DeleteAllEntriesForKeyPrefix(ctx sdk.Context, keyPrefix []byte) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, keyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// ---------------------- Config Stores  ---------------------- //

// GetDaysSinceModuleGenesis returns the number of days since the module was initialized
func (k Keeper) GetDaysSinceModuleGenesis(ctx sdk.Context) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDaysSinceGenesis)
	bz := store.Get(types.KeyPrefixDaysSinceGenesis)
	if bz == nil {
		// This should never happen as the module is initialized with 0 days on genesis
		return 0, fmt.Errorf("days since module genesis not found")
	}

	daysSinceGenesis := sdk.BigEndianToUint64(bz)

	return daysSinceGenesis, nil
}

// SetDaysSinceModuleGenesis updates the number of days since genesis
func (k Keeper) SetDaysSinceModuleGenesis(ctx sdk.Context, daysSinceGenesis uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDaysSinceGenesis)
	store.Set(types.KeyPrefixDaysSinceGenesis, sdk.Uint64ToBigEndian(daysSinceGenesis))
}

// GetDeveloperFees returns the fees the developers can withdraw from the module account
func (k Keeper) GetDeveloperFees(ctx sdk.Context, denom string) (sdk.Coin, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeveloperFees)
	key := types.GetKeyPrefixDeveloperFees(denom)

	bz := store.Get(key)
	if bz == nil {
		return sdk.Coin{}, fmt.Errorf("developer fees for %s not found", denom)
	}

	developerFees := sdk.Coin{}
	err := developerFees.Unmarshal(bz)
	if err != nil {
		return sdk.Coin{}, err
	}

	return developerFees, nil
}

// GetAllDeveloperFees returns all the developer fees the developer account can withdraw
func (k Keeper) GetAllDeveloperFees(ctx sdk.Context) ([]sdk.Coin, error) {
	fees := make([]sdk.Coin, 0)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeveloperFees)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixDeveloperFees)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		developerFees := sdk.Coin{}
		if err := developerFees.Unmarshal(iterator.Value()); err != nil {
			return nil, fmt.Errorf("error unmarshalling developer fees: %w", err)
		}

		fees = append(fees, developerFees)
	}

	return fees, nil
}

// SetDeveloperFees sets the fees the developers can withdraw from the module account
func (k Keeper) SetDeveloperFees(ctx sdk.Context, developerFees sdk.Coin) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeveloperFees)
	key := types.GetKeyPrefixDeveloperFees(developerFees.Denom)

	bz, err := developerFees.Marshal()
	if err != nil {
		return err
	}

	store.Set(key, bz)

	return nil
}

// DeleteDeveloperFees deletes the developer fees given a denom
func (k Keeper) DeleteDeveloperFees(ctx sdk.Context, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeveloperFees)
	key := types.GetKeyPrefixDeveloperFees(denom)
	store.Delete(key)
}

// GetProtoRevEnabled returns whether protorev is enabled
func (k Keeper) GetProtoRevEnabled(ctx sdk.Context) (bool, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtoRevEnabled)
	bz := store.Get(types.KeyPrefixProtoRevEnabled)
	if bz == nil {
		// This should never happen as the module is initialized on genesis
		return false, fmt.Errorf("protorev enabled/disabled configuration has not been set in state")
	}

	res, err := strconv.ParseBool(string(bz))
	if err != nil {
		return false, err
	}

	return res, nil
}

// SetProtoRevEnabled sets whether the protorev post handler is enabled
func (k Keeper) SetProtoRevEnabled(ctx sdk.Context, enabled bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixProtoRevEnabled)
	bz := []byte(strconv.FormatBool(enabled))
	store.Set(types.KeyPrefixProtoRevEnabled, bz)
}

// GetPointCountForBlock returns the number of pool points that have been consumed in the current block
func (k Keeper) GetPointCountForBlock(ctx sdk.Context) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPointCountForBlock)
	bz := store.Get(types.KeyPrefixPointCountForBlock)
	if bz == nil {
		// This should never happen as this is set to 0 on genesis
		return 0, fmt.Errorf("current pool point count has not been set in state")
	}

	res := sdk.BigEndianToUint64(bz)

	return res, nil
}

// SetPointCountForBlock sets the number of pool points that have been consumed in the current block
func (k Keeper) SetPointCountForBlock(ctx sdk.Context, txCount uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPointCountForBlock)
	store.Set(types.KeyPrefixPointCountForBlock, sdk.Uint64ToBigEndian(txCount))
}

// IncrementPointCountForBlock increments the number of pool points that have been consumed in the current block
func (k Keeper) IncrementPointCountForBlock(ctx sdk.Context, amount uint64) error {
	pointCount, err := k.GetPointCountForBlock(ctx)
	if err != nil {
		return err
	}

	k.SetPointCountForBlock(ctx, pointCount+amount)

	return nil
}

// GetLatestBlockHeight returns the latest block height that protorev was run on
func (k Keeper) GetLatestBlockHeight(ctx sdk.Context) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLatestBlockHeight)
	bz := store.Get(types.KeyPrefixLatestBlockHeight)
	if bz == nil {
		// This should never happen as the module is initialized on genesis and reset in the post handler
		return 0, fmt.Errorf("block height has not been set in state")
	}

	res := sdk.BigEndianToUint64(bz)

	return res, nil
}

// SetLatestBlockHeight sets the latest block height that protorev was run on
func (k Keeper) SetLatestBlockHeight(ctx sdk.Context, blockHeight uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLatestBlockHeight)
	store.Set(types.KeyPrefixLatestBlockHeight, sdk.Uint64ToBigEndian(blockHeight))
}

// ---------------------- Admin Stores  ---------------------- //

// GetAdminAccount returns the admin account for protorev
func (k Keeper) GetAdminAccount(ctx sdk.Context) (sdk.AccAddress, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAdminAccount)
	bz := store.Get(types.KeyPrefixAdminAccount)
	if bz == nil {
		return nil, fmt.Errorf("admin account not found, it has not been initialized through governance")
	}

	return sdk.AccAddress(bz), nil
}

// SetAdminAccount sets the admin account for protorev
func (k Keeper) SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAdminAccount)
	store.Set(types.KeyPrefixAdminAccount, adminAccount.Bytes())
}

// GetDeveloperAccount returns the developer account for protorev
func (k Keeper) GetDeveloperAccount(ctx sdk.Context) (sdk.AccAddress, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeveloperAccount)
	bz := store.Get(types.KeyPrefixDeveloperAccount)
	if bz == nil {
		return nil, fmt.Errorf("developer account not found, it has not been initialized by the admin account")
	}

	return sdk.AccAddress(bz), nil
}

// SetDeveloperAccount sets the developer account for protorev that will receive a portion of arbitrage profits
func (k Keeper) SetDeveloperAccount(ctx sdk.Context, developerAccount sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDeveloperAccount)
	store.Set(types.KeyPrefixDeveloperAccount, developerAccount.Bytes())
}

// GetMaxPointsPerTx returns the max number of pool points that can be consumed per transaction. A pool point is roughly
// equivalent to 1 ms of simulation & execution time.
func (k Keeper) GetMaxPointsPerTx(ctx sdk.Context) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixMaxPointsPerTx)
	bz := store.Get(types.KeyPrefixMaxPointsPerTx)
	if bz == nil {
		// This should never happen as it is set to the default value on genesis
		return 0, fmt.Errorf("max pool points per tx has not been set in state")
	}

	res := sdk.BigEndianToUint64(bz)
	return res, nil
}

// SetMaxPointsPerTx sets the max number of pool points that can be consumed per transaction. A pool point is roughly
// equivalent to 1 ms of simulation & execution time.
func (k Keeper) SetMaxPointsPerTx(ctx sdk.Context, maxPoints uint64) error {
	if maxPoints == 0 || maxPoints > types.MaxPoolPointsPerTx {
		return fmt.Errorf("max pool points must be between 1 and %d", types.MaxPoolPointsPerTx)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixMaxPointsPerTx)
	bz := sdk.Uint64ToBigEndian(maxPoints)
	store.Set(types.KeyPrefixMaxPointsPerTx, bz)

	return nil
}

// GetMaxPointsPerBlock returns the max number of pool points that can be consumed per block. A pool point is roughly
// equivalent to 1 ms of simulation & execution time.
func (k Keeper) GetMaxPointsPerBlock(ctx sdk.Context) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixMaxPointsPerBlock)
	bz := store.Get(types.KeyPrefixMaxPointsPerBlock)
	if bz == nil {
		// This should never happen as it is set to the default value on genesis
		return 0, fmt.Errorf("max pool points per block has not been set in state")
	}

	res := sdk.BigEndianToUint64(bz)
	return res, nil
}

// SetMaxPointsPerBlock sets the max number of pool points that can be consumed per block. A pool point is roughly
// equivalent to 1 ms of simulation & execution time.
func (k Keeper) SetMaxPointsPerBlock(ctx sdk.Context, maxPoints uint64) error {
	if maxPoints == 0 || maxPoints > types.MaxPoolPointsPerBlock {
		return fmt.Errorf("max pool points per block must be between 1 and %d", types.MaxPoolPointsPerBlock)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixMaxPointsPerBlock)
	bz := sdk.Uint64ToBigEndian(maxPoints)
	store.Set(types.KeyPrefixMaxPointsPerBlock, bz)

	return nil
}

// GetPoolWeights retrieves the weights of different pool types. The weight of a pool type roughly
// corresponds to the amount of time it will take to simulate and execute a swap on that pool type (in ms).
func (k Keeper) GetPoolWeights(ctx sdk.Context) *types.PoolWeights {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPoolWeights)
	poolWeights := &types.PoolWeights{}
	osmoutils.MustGet(store, types.KeyPrefixPoolWeights, poolWeights)
	return poolWeights
}

// SetPoolWeights sets the weights of different pool types.
func (k Keeper) SetPoolWeights(ctx sdk.Context, poolWeights types.PoolWeights) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPoolWeights)
	osmoutils.MustSet(store, types.KeyPrefixPoolWeights, &poolWeights)
}
