package types

import (
	"strings"

	"github.com/dop251/goja"

	"github.com/smartbch/egvm/script/utils"
)

type BufBuilder struct {
	sb *strings.Builder
}

func (b BufBuilder) Len() int {
	return b.sb.Len()
}

func (b BufBuilder) Reset() {
	b.sb.Reset()
}

func NewBufBuilder() BufBuilder {
	sb := &strings.Builder{}
	return BufBuilder{sb: sb}
}

func (b BufBuilder) Write(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	bz := utils.GetOneArrayBuffer(f)
	n, err := b.sb.Write(bz)
	if err != nil {
		panic(goja.NewSymbol("error in Ecrecover: " + err.Error()))
	}
	return vm.ToValue(n)
}

func (b BufBuilder) ToBuf(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(utils.IncorrectArgumentCount)
	}
	return vm.ToValue(vm.NewArrayBuffer([]byte(b.sb.String())))
}
