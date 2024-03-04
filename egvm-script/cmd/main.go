package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"log"

	"github.com/dop251/goja"
)

func main() {
	var randomSeedStr string
	var fileName string
	flag.StringVar(&randomSeedStr, "r", "hello-world", "Random Seed String")
	flag.StringVar(&fileName, "f", "script.js", "Script File")
	flag.Parse()
	script, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	vm := goja.New()
	registerFunctions(vm)
	_, err = vm.RunString(string(script))
	if err != nil {
		log.Fatalf("unable to execute script file: %v", err)
	}

	setRandomSeed, ok := goja.AssertFunction(vm.Get("setRandomSeed"))
	if !ok {
		log.Fatalf("cannot find the 'setRandomSeed' function")
	}
	handleInput, ok := goja.AssertFunction(vm.Get("handleInput"))
	if !ok {
		log.Fatalf("cannot find the 'handleInput' function")
	}

	_, err = setRandomSeed(goja.Undefined(), vm.ToValue(randomSeedStr))
	if err != nil {
		log.Fatalf("failed to run setRandomSeed: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		res, _ := handleInput(goja.Undefined(), vm.ToValue(scanner.Text()))
		var outStr string
		vm.ExportTo(res, &outStr)
		fmt.Println(outStr)
	}
}

