package indexer

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/light"
	tmtypes "github.com/tendermint/tendermint/types"
)

const RecoveryFileEveryNBlocks = 1000
const TendermintChainID = "0x2710"

// loop get block and events，verify block
// register contract jobs
type DataFetcher struct {
	tmServerAddrList    []string
	latestTrustedHeader *tmtypes.SignedHeader
	contractAddress     common.Address
	ctx                 *Context
}

// loop get blocks
// todo：loop get events
func (df *DataFetcher) start() {
	for {
		height := df.latestTrustedHeader.Height + 1
		untrustedHeader := getSignedHeader(df.tmServerAddrList, uint64(height))
		if untrustedHeader.Height != height {
			df.tmServerAddrList = rotateList(df.tmServerAddrList)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		validators := getValidators(df.tmServerAddrList, uint64(height))
		valSet, err := tmtypes.ValidatorSetFromExistingValidators(validators)
		lb := &tmtypes.LightBlock{
			SignedHeader: untrustedHeader,
			ValidatorSet: valSet,
		}
		err = lb.ValidateBasic(TendermintChainID)
		if err != nil {
			fmt.Printf("light block at height [%d] verify failed\n", height)
			df.tmServerAddrList = rotateList(df.tmServerAddrList)
			continue
		}
		err = light.VerifyAdjacent(df.latestTrustedHeader, untrustedHeader, valSet, 2*168*time.Hour, time.Now(), 10*time.Second)
		if err != nil {
			fmt.Printf("header at height [%d] verify failed\n", height)
			df.tmServerAddrList = rotateList(df.tmServerAddrList)
			continue
		}
		blk := getBlock(df.tmServerAddrList, uint64(height))
		if bH, tH := blk.Hash(), lb.Hash(); !bytes.Equal(bH, tH) {
			fmt.Printf("block hash not equal signedHeader hash at height [%d]\n", height)
			df.tmServerAddrList = rotateList(df.tmServerAddrList)
			continue
		}
		df.latestTrustedHeader = untrustedHeader
		df.ctx.HandleBlock(df.contractAddress, blk) 
		if blk.Height % RecoveryFileEveryNBlocks == 0 {
			recoveryFileName := fmt.Sprintf("%09d.dump", blk.Height)
			df.ctx.DumpRecoveryFile(recoveryFileName, df.latestTrustedHeader)
		}
	}
}

// if something wrong with current chain node， not discard it， just rotate another one front，log it
func rotateList(smartbchAddrs []string) []string {
	var newList = []string{smartbchAddrs[len(smartbchAddrs)-1]}
	newList = append(newList, smartbchAddrs[:len(smartbchAddrs)-1]...)
	return newList
}
