package main

import (
	"bytes"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// recovery job from db when program restart
type Job struct {
	indexer *Indexer

	address     string
	beginHeight int64
}

type TxData struct {
	From  string
	Value string
	Input []byte
}

// scan blocks and events from beginHeight
func (j *Job) Start() {
	indexer := j.indexer
	h := j.beginHeight
	for {
		txs := indexer.GetTxsByHeight(h)
		if txs == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		var datas []TxData
		for _, txRawData := range txs {
			tx := &types.Transaction{}
			err := tx.DecodeRLP(rlp.NewStream(bytes.NewReader(txRawData), 0))
			if err != nil {
				continue
			}
			if tx.To().String() == j.address {
				from, err := indexer.signer.Sender(tx)
				if err != nil {
					continue
				}
				data := TxData{
					From:  from.String(),
					Value: tx.Value().String(),
					Input: tx.Data(),
				}
				datas = append(datas, data)
			}
		}
		if datas != nil {
			indexer.storeTxDatasByAddressAndHeight(j.address, h, datas)
		}
		h++
	}
}
