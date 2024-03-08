package utils

import "github.com/dop251/goja"

func GetOneUint64(f goja.FunctionCall) uint64 {
	if len(f.Arguments) != 1 {
		panic(IncorrectArgumentCount)
	}
	a, ok := f.Arguments[0].Export().(int64)
	if !ok {
		panic(goja.NewSymbol("The first argument must be number"))
	}

	if uint64(a) > MaxSafeInteger {
		panic(LargerThanMaxInteger)
	}

	return uint64(a)
}

func GetOneArrayBuffer(f goja.FunctionCall) []byte {
	if len(f.Arguments) != 1 {
		panic(IncorrectArgumentCount)
	}
	a, ok := f.Arguments[0].Export().(goja.ArrayBuffer)
	if !ok {
		panic(goja.NewSymbol("The first argument must be ArrayBuffer"))
	}
	return a.Bytes()
}

func GetTwoArrayBuffers(f goja.FunctionCall) ([]byte, []byte) {
	if len(f.Arguments) != 2 {
		panic(IncorrectArgumentCount)
	}
	a, ok := f.Arguments[0].Export().(goja.ArrayBuffer)
	if !ok {
		panic(goja.NewSymbol("The first argument must be ArrayBuffer"))
	}
	b, ok := f.Arguments[1].Export().(goja.ArrayBuffer)
	if !ok {
		panic(goja.NewSymbol("The second argument must be ArrayBuffer"))
	}
	return a.Bytes(), b.Bytes()
}

func GetThreeArrayBuffers(f goja.FunctionCall) ([]byte, []byte, []byte) {
	if len(f.Arguments) != 3 {
		panic(IncorrectArgumentCount)
	}
	a, ok := f.Arguments[0].Export().(goja.ArrayBuffer)
	if !ok {
		panic(goja.NewSymbol("The first argument must be ArrayBuffer"))
	}
	b, ok := f.Arguments[1].Export().(goja.ArrayBuffer)
	if !ok {
		panic(goja.NewSymbol("The second argument must be ArrayBuffer"))
	}
	c, ok := f.Arguments[2].Export().(goja.ArrayBuffer)
	if !ok {
		panic(goja.NewSymbol("The second argument must be ArrayBuffer"))
	}
	return a.Bytes(), b.Bytes(), c.Bytes()
}
