package utils

import (
	"github.com/dop251/goja"
	"github.com/holiman/uint256"
)

const (
	MaxSafeInteger = (uint64(1) << 53) - 1
)

var (
	MinNegValue = uint256.NewInt(0).Lsh(uint256.NewInt(1), 255)
)

// error
var (
	// function error
	IncorrectArgumentCount = goja.NewSymbol("Incorrect argument count")

	// ordered map error
	EmptyKeyString = goja.NewSymbol("Empty key string")

	// number overflow
	LargerThanMaxInteger = goja.NewSymbol("Larger than Number.MAX_SAFE_INTEGER")
	OverflowInSigned     = goja.NewSymbol("Overflow in signed multiplication")
)
