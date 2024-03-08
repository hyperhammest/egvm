package indexer

import (
	"crypto/sha256"
	"encoding/json"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartbch/egvm/keygrantor"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	KeyFile      = "/data/key.txt"
	KGClient *keygrantor.SimpleClient
)

func start() (df DataFetcher) {
	var scriptFileName string
	var contractAddress string
	var headerFileName string
	var keyGrantorUrl string
	var recoveryFileName string
	var tmServerListArg string
	flag.StringVar(&scriptFileName, "script", "script.js", "Script File")
	flag.StringVar(&contractAddress, "contract", "", "The Game's Smart Contract Address")
	flag.StringVar(&headerFileName, "header", "header.json", "Header File")
	flag.StringVar(&keyGrantorUrl, "keygrantor", "", "KeyGrantor's URL")
	flag.StringVar(&recoveryFileName, "recover", "", "Recovery File")
	flag.StringVar(&tmServerListArg, "tmservers", "", "Servers Providing Tendermint RPC")
	flag.Parse()

	df.tmServerAddrList = strings.Split(tmServerListArg, ",")

	script, err := ioutil.ReadFile(scriptFileName)
	if err != nil {
		log.Fatalf("unable to read script file: %v", err)
	}
	scriptHash := sha256.Sum256(script)
	vm := goja.New()
	registerFunctions(vm)
	_, err = vm.RunString(string(script))
	if err != nil {
		log.Fatalf("unable to execute script file: %v", err)
	}

	if !common.IsHexAddress(contractAddress) {
		log.Fatalf("invalid address: %s", contractAddress)
	}
	df.contractAddress = common.HexToAddress(contractAddress)

	header, err := ioutil.ReadFile(headerFileName)
	if err != nil {
		log.Fatalf("unable to read header file: %v", err)
	}
	headerHash := sha256.Sum256(header)
	var trustedHeader   tmtypes.SignedHeader
	err = json.Unmarshal(header, &trustedHeader)
	if err != nil {
		log.Fatalf("unable to unmarshal header: %v", err)
	}
	df.latestTrustedHeader = &trustedHeader

	configHash := sha256.Sum256(append(scriptHash[:], append(headerHash[:], df.contractAddress[:]...)...))
	_, err = os.ReadFile(KeyFile)
	if err != nil {
		if os.IsNotExist(err) {
			KGClient.InitKeys(keyGrantorUrl, configHash, false)
			fmt.Printf("get enclave vrf private key from keygrantor, its pubkey is: %s\n", hex.EncodeToString(KGClient.PubKeyBz))
			keygrantor.SealKeyToFile(KeyFile, KGClient.ExtPrivKey)
		} else {
			panic(err)
		}
	} else {
		KGClient.InitKeys(KeyFile, configHash, true)
	}

	if len(recoveryFileName) == 0 {
		df.ctx = NewContext(vm)
	} else {
		df.ctx = NewContextFromRecoverFile(recoveryFileName, vm)
		df.latestTrustedHeader = df.ctx.TrustedHeader
	}
	return
}

