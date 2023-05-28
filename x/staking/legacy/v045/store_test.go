package v045_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v045staking "github.com/cosmos/cosmos-sdk/x/staking/legacy/v045"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestMigrate(t *testing.T) {
	encCfg := simapp.MakeTestEncodingConfig()
	stakingKey := sdk.NewKVStoreKey("staking")
	tStakingKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(stakingKey, tStakingKey)
	paramstore := paramtypes.NewSubspace(encCfg.Marshaler, encCfg.Amino, stakingKey, tStakingKey, "staking")

	testCases := []struct {
		OldValidator types.Validator
		NewValidator types.Validator
	}{
		{
			OldValidator: types.Validator{
				OperatorAddress: sdk.ValAddress{0x00}.String(),
				Commission: types.Commission{
					CommissionRates: types.CommissionRates{
						Rate:    sdk.MustNewDecFromStr("0.01"),
						MaxRate: sdk.MustNewDecFromStr("0.02"),
					},
				},
			},
			NewValidator: types.Validator{
				OperatorAddress: sdk.ValAddress{0x00}.String(),
				Commission: types.Commission{
					CommissionRates: types.CommissionRates{
						Rate:    sdk.MustNewDecFromStr("0.05"),
						MaxRate: sdk.MustNewDecFromStr("0.05"),
					},
				},
			},
		},
		{
			OldValidator: types.Validator{
				OperatorAddress: sdk.ValAddress{0x01}.String(),
				Commission: types.Commission{
					CommissionRates: types.CommissionRates{
						Rate:    sdk.MustNewDecFromStr("0.05"),
						MaxRate: sdk.MustNewDecFromStr("0.05"),
					},
				},
			},
			NewValidator: types.Validator{
				OperatorAddress: sdk.ValAddress{0x01}.String(),
				Commission: types.Commission{
					CommissionRates: types.CommissionRates{
						Rate:    sdk.MustNewDecFromStr("0.05"),
						MaxRate: sdk.MustNewDecFromStr("0.05"),
					},
				},
			},
		},
		{
			OldValidator: types.Validator{
				OperatorAddress: sdk.ValAddress{0x02}.String(),
				Commission: types.Commission{
					CommissionRates: types.CommissionRates{
						Rate:    sdk.MustNewDecFromStr("0.1"),
						MaxRate: sdk.MustNewDecFromStr("0.2"),
					},
				},
			},
			NewValidator: types.Validator{
				OperatorAddress: sdk.ValAddress{0x02}.String(),
				Commission: types.Commission{
					CommissionRates: types.CommissionRates{
						Rate:    sdk.MustNewDecFromStr("0.1"),
						MaxRate: sdk.MustNewDecFromStr("0.2"),
					},
				},
			},
		},
	}

	store := ctx.KVStore(stakingKey)
	for _, vs := range testCases {
		bz := types.MustMarshalValidator(encCfg.Marshaler, &vs.OldValidator)
		store.Set(types.GetValidatorKey(vs.OldValidator.GetOperator()), bz)
	}

	// Check no params
	require.False(t, paramstore.Has(ctx, types.KeyMinCommissionRate))

	// Run migrations.
	v045staking.MigrateStore(ctx, stakingKey, encCfg.Marshaler, paramstore)

	// Make sure the new params are set.
	require.True(t, paramstore.Has(ctx, types.KeyMinCommissionRate))

	minCommissionRate := sdk.Dec{}
	paramstore.Get(ctx, types.KeyMinCommissionRate, &minCommissionRate)
	require.Equal(t, types.DefaultMinCommissionRate, minCommissionRate)

	// Make sure the validators commission is set correctly.
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(encCfg.Marshaler, iterator.Value())
		require.Equal(t, testCases[i].NewValidator.OperatorAddress, validator.GetOperator().String())
		require.Equal(t, testCases[i].NewValidator.Commission.Rate, validator.Commission.Rate)
		require.Equal(t, testCases[i].NewValidator.Commission.MaxRate, validator.Commission.MaxRate)
		i++
	}
}
