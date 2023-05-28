package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Migrate performs in-place store migrations from v0.45.13 to v0.45.14.
// The migration includes:
//
// - Adding MinCommissionRate param
// - Setting validaotr commission rate and max commission rate to MinCommissionRate if they are lower
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, paramstore paramtypes.Subspace) {
	migrateParamsStore(ctx, paramstore)
	migrateValidators(ctx, storeKey, cdc)
}

func migrateParamsStore(ctx sdk.Context, paramstore paramtypes.Subspace) {
	if paramstore.HasKeyTable() {
		paramstore.Set(ctx, types.KeyMinCommissionRate, types.DefaultMinCommissionRate)
	} else {
		paramstore.WithKeyTable(types.ParamKeyTable())
		paramstore.Set(ctx, types.KeyMinCommissionRate, types.DefaultMinCommissionRate)
	}
}

func migrateValidators(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(cdc, iterator.Value())
		if validator.Commission.CommissionRates.Rate.LT(types.DefaultMinCommissionRate) {
			validator.Commission.CommissionRates.Rate = types.DefaultMinCommissionRate
		}

		if validator.Commission.CommissionRates.MaxRate.LT(types.DefaultMinCommissionRate) {
			validator.Commission.CommissionRates.MaxRate = types.DefaultMinCommissionRate
		}
		store.Set(iterator.Key(), types.MustMarshalValidator(cdc, &validator))
	}
}
