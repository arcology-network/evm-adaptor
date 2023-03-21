package evm

import "math/big"

type txContext struct {
	height *big.Int
	index  uint32
}

func (c *txContext) GetIndex() uint32 {
	return c.index
}

func (c *txContext) GetHeight() *big.Int {
	return c.height
}
