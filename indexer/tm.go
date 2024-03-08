package indexer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	tmjson "github.com/tendermint/tendermint/libs/json"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func getSignedHeader(addrs []string, height uint64) *tmtypes.SignedHeader {
	for {
		for _, addr := range addrs {
			if !strings.Contains(addr, "8545") {
				continue
			}
			addr = strings.ReplaceAll(addr, "8545", "26657")
			b, err := getSignedHeaderByNumber(addr, height)
			if err == nil {
				return b
			}
		}
		fmt.Println("retry getSignedHeader")
		time.Sleep(10 * time.Second)
	}
}

func getSignedHeaderByNumber(url string, height uint64) (*tmtypes.SignedHeader, error) {
	reqUrl := fmt.Sprintf("%s/commit?height=%d", url, height)
	resp, err := http.Get(reqUrl)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return nil, errors.New(fmt.Sprintf("response status not wanted:%s", resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var r rpctypes.RPCResponse
	err = r.UnmarshalJSON(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var c ctypes.ResultCommit
	err = tmjson.Unmarshal(r.Result, &c)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &c.SignedHeader, nil
}

func getValidators(addrs []string, height uint64) []*tmtypes.Validator {
	for {
		for _, addr := range addrs {
			if !strings.Contains(addr, "8545") {
				continue
			}
			addr = strings.ReplaceAll(addr, "8545", "26657")
			validators, err := getValidatorsByNumber(addr, height)
			if err == nil {
				return validators
			}
		}
		fmt.Println("retry getValidators")
		time.Sleep(10 * time.Second)
	}
}

func getValidatorsByNumber(url string, height uint64) ([]*tmtypes.Validator, error) {
	reqUrl := fmt.Sprintf("%s/validators?height=%d", url, height)
	resp, err := http.Get(reqUrl)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return nil, errors.New(fmt.Sprintf("response status not wanted:%s", resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var r rpctypes.RPCResponse
	err = r.UnmarshalJSON(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var b ctypes.ResultValidators
	err = tmjson.Unmarshal(r.Result, &b)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return b.Validators, nil
}

func getBlock(addrs []string, height uint64) *tmtypes.Block {
	for {
		for _, addr := range addrs {
			if !strings.Contains(addr, "8545") {
				continue
			}
			addr = strings.ReplaceAll(addr, "8545", "26657")
			b, err := getBlockByNumber(addr, height)
			if err == nil {
				return b
			}
		}
		fmt.Println("retry getBlock")
		time.Sleep(10 * time.Second)
	}
}

func getBlockByNumber(url string, height uint64) (*tmtypes.Block, error) {
	reqUrl := fmt.Sprintf("%s/block?height=%d", url, height)
	resp, err := http.Get(reqUrl)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return nil, errors.New(fmt.Sprintf("response status not wanted:%s", resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var r rpctypes.RPCResponse
	err = r.UnmarshalJSON(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var b ctypes.ResultBlock
	err = tmjson.Unmarshal(r.Result, &b)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return b.Block, nil
}
