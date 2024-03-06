package main

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	latestTrustedHeaderKey = []byte("egvm-header")
)

func buildTxsKey(height int64) []byte {
	return []byte(fmt.Sprintf("egvm-txs-%d", height))
}

func buildTxDatasKey(address string, height int64) []byte {
	return []byte(fmt.Sprintf("egvm-data-%s-%d", address, height))
}

func (i *Indexer) storeTxDatasByAddressAndHeight(address string, height int64, datas []TxData) {
	err := i.db.Update(func(txn *badger.Txn) error {
		bz, _ := json.Marshal(datas)
		err := txn.Set(buildTxDatasKey(address, height), bz)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func (i *Indexer) getTxDatasByAddressAndHeight(address string, height int64) (datas []TxData) {
	err := i.db.View(func(txn *badger.Txn) error {
		iterm, err := txn.Get(buildTxDatasKey(address, height))
		if err != nil {
			return err
		}
		return iterm.Value(func(val []byte) error {
			var bz = make([]byte, len(val))
			copy(bz, val)
			err := json.Unmarshal(bz, &datas)
			if err != nil {
				return err
			}
			return nil
		})
	})
	if err == badger.ErrKeyNotFound {
		return nil
	} else {
		panic(err)
	}
}

func (i *Indexer) storeTxsByHeight(height int64, txs []tmtypes.Tx) {
	err := i.db.Update(func(txn *badger.Txn) error {
		bz, _ := json.Marshal(txs)
		err := txn.Set(buildTxsKey(height), bz)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func (i *Indexer) GetTxsByHeight(height int64) (txs []tmtypes.Tx) {
	err := i.db.View(func(txn *badger.Txn) error {
		iterm, err := txn.Get(buildTxsKey(height))
		if err != nil {
			return err
		}
		return iterm.Value(func(val []byte) error {
			var bz = make([]byte, len(val))
			copy(bz, val)
			err := json.Unmarshal(bz, &txs)
			if err != nil {
				return err
			}
			return nil
		})
	})
	if err == badger.ErrKeyNotFound {
		return nil
	} else {
		panic(err)
	}
}

func (i *Indexer) storeLatestTrustedHeader(h tmtypes.SignedHeader) {
	err := i.db.Update(func(txn *badger.Txn) error {
		bz, _ := tmjson.Marshal(h)
		err := txn.Set(latestTrustedHeaderKey, bz)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func (i *Indexer) recoveryLatestTrustedHeader() error {
	return i.db.View(func(txn *badger.Txn) error {
		iterm, err := txn.Get(latestTrustedHeaderKey)
		if err != nil {
			return err
		}
		return iterm.Value(func(val []byte) error {
			var bz = make([]byte, len(val))
			copy(bz, val)
			var h tmtypes.SignedHeader
			err := tmjson.Unmarshal(bz, &h)
			if err != nil {
				return err
			}
			i.latestTrustedHeader = &h
			return nil
		})
	})
}
