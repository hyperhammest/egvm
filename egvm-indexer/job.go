package main

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type Job struct {
	indexer *Indexer

	address     string
	beginHeight int64
}

func (j *Job) Start() {
	// todo: get txs from db since height job.beginHeight
	var txs [][]byte
	for _, txData := range txs {
		tx := &types.Transaction{}
		err := tx.DecodeRLP(rlp.NewStream(bytes.NewReader(txData), 0))
		if err != nil {
			continue
		}
		if tx.To().String() == j.address {
			// todo: store tx.Data()
		}
	}
}
