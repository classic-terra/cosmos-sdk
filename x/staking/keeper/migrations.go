package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v043 "github.com/cosmos/cosmos-sdk/x/staking/legacy/v043"
	v045 "github.com/cosmos/cosmos-sdk/x/staking/legacy/v045"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v043.MigrateStore(ctx, m.keeper.storeKey)
}

// Migrate13to14 migrates from version v0.45.13 to v0.45.14.
// Only for this particular version, which do not use the version of module.
func (m Migrator) Migrate13to14(ctx sdk.Context) {
	v045.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.keeper.paramstore)
}
