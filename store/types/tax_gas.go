package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

// Gas measured by the SDK
type TaxGas = sdkmath.Uint

// GasMeter interface to track gas consumption
type TaxGasMeter interface {
	GasConsumed() TaxGas
	ConsumeGas(amount TaxGas, descriptor string)
	RefundGas(amount TaxGas, descriptor string)
	String() string
}

type basicTaxGasMeter struct {
	consumed TaxGas
}

// NewGasMeter returns a reference to a new basicGasMeter.
func NewTaxGasMeter() TaxGasMeter {
	return &basicTaxGasMeter{
		consumed: sdkmath.ZeroUint(),
	}
}

// GasConsumed returns the gas consumed from the GasMeter.
func (g *basicTaxGasMeter) GasConsumed() TaxGas {
	return g.consumed
}

// ConsumeGas adds the given amount of gas to the gas consumed and panics if it overflows the limit or out of gas.
func (g *basicTaxGasMeter) ConsumeGas(amount TaxGas, _descriptor string) {
	g.consumed = g.consumed.Add(amount)
}

// RefundGas will deduct the given amount from the gas consumed. If the amount is greater than the
// gas consumed, the function will panic.
//
// Use case: This functionality enables refunding gas to the transaction or block gas pools so that
// EVM-compatible chains can fully support the go-ethereum StateDb interface.
// See https://github.com/cosmos/cosmos-sdk/pull/9403 for reference.
func (g *basicTaxGasMeter) RefundGas(amount TaxGas, descriptor string) {
	if g.consumed.LTE(amount) {
		panic(ErrorNegativeGasConsumed{Descriptor: descriptor})
	}

	g.consumed = g.consumed.Sub(amount)
}

// String returns the BasicGasMeter's gas limit and gas consumed.
func (g *basicTaxGasMeter) String() string {
	return fmt.Sprintf("BasicTaxGasMeter:\n  consumed: %s", g.consumed.String())
}
