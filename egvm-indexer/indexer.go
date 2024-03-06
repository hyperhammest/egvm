package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/light"
	tmtypes "github.com/tendermint/tendermint/types"
)

// loop get block and events，verify block and store it
// register contract jobs
type Indexer struct {
	chainId          string
	signer           gethtypes.Signer
	smartBCHAddrList []string
	listenUrl        string

	latestTrustedHeader *tmtypes.SignedHeader
	db                  *badger.DB

	jobList map[string]*Job
}

// loop get blocks
// todo：loop get events
func (i *Indexer) start() {
	for {
		height := i.latestTrustedHeader.Height + 1
		untrustedHeader := getSignedHeader(i.smartBCHAddrList, uint64(height))
		if untrustedHeader.Height != height {
			i.smartBCHAddrList = rotateList(i.smartBCHAddrList)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		validators := getValidators(i.smartBCHAddrList, uint64(height))
		valSet, err := tmtypes.ValidatorSetFromExistingValidators(validators)
		lb := &tmtypes.LightBlock{
			SignedHeader: untrustedHeader,
			ValidatorSet: valSet,
		}
		err = lb.ValidateBasic(i.chainId)
		if err != nil {
			fmt.Printf("light block at height [%d] verify failed\n", height)
			i.smartBCHAddrList = rotateList(i.smartBCHAddrList)
			continue
		}
		err = light.VerifyAdjacent(i.latestTrustedHeader, untrustedHeader, valSet, 2*168*time.Hour, time.Now(), 10*time.Second)
		if err != nil {
			fmt.Printf("header at height [%d] verify failed\n", height)
			i.smartBCHAddrList = rotateList(i.smartBCHAddrList)
			continue
		}
		blk := getBlock(i.smartBCHAddrList, uint64(height))
		if bH, tH := blk.Hash(), lb.Hash(); !bytes.Equal(bH, tH) {
			fmt.Printf("block hash not equal signedHeader hash at height [%d]\n", height)
			i.smartBCHAddrList = rotateList(i.smartBCHAddrList)
			continue
		}
		i.storeTxsByHeight(height, blk.Txs)
		i.latestTrustedHeader = untrustedHeader
		i.storeLatestTrustedHeader(*i.latestTrustedHeader)
	}
}

// todo： recovery all jobs
// todo： move initLatestTrustedHeader here
func (i *Indexer) init() {

}

func (i *Indexer) initLatestTrustedHeader() {
	err := i.recoveryLatestTrustedHeader()
	if err == badger.ErrKeyNotFound {
		var h tmtypes.SignedHeader
		err = tmjson.Unmarshal([]byte(InitialHeaderStr), &h)
		if err != nil {
			panic(err)
		}
		i.latestTrustedHeader = &h
	} else {
		panic(err)
	}
}

// add new contract index job for https endpoint
func (i *Indexer) AddJob(address string, height int64) error {
	if _, exist := i.jobList[address]; exist {
		return jobExistErr
	}
	if len(i.jobList) > maxJobNums {
		return tooManyJobsErr
	}
	j := &Job{indexer: i, beginHeight: height}
	i.jobList[address] = j
	j.Start()
	return nil
}

// if something wrong with current chain node， not discard it， just rotate another one front，log it
func rotateList(smartbchAddrs []string) []string {
	var newList = []string{smartbchAddrs[len(smartbchAddrs)-1]}
	newList = append(newList, smartbchAddrs[:len(smartbchAddrs)-1]...)
	return newList
}
