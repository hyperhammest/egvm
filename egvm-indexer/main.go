package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/light"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	jobExistErr    = errors.New("job already exist")
	tooManyJobsErr = errors.New("job nums touch upper limit")
)

var (
	latestTrustedHeaderKey = []byte("egvm-header")
)

const (
	maxJobNums = 1000
)

type Indexer struct {
	chainId          string
	smartBCHAddrList []string
	listenUrl        string

	latestTrustedHeader *tmtypes.SignedHeader
	db                  *badger.DB

	jobList map[string]*Job
}

func (i *Indexer) start() {
	for {
		height := i.latestTrustedHeader.Height
		untrustedHeader := getSignedHeader(i.smartBCHAddrList, uint64(height+1))
		if untrustedHeader.Height != height+1 {
			i.smartBCHAddrList = rotateList(i.smartBCHAddrList)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		validators := getValidators(i.smartBCHAddrList, uint64(height+1))
		valSet, err := tmtypes.ValidatorSetFromExistingValidators(validators)
		lb := &tmtypes.LightBlock{
			SignedHeader: untrustedHeader,
			ValidatorSet: valSet,
		}
		err = lb.ValidateBasic(i.chainId)
		if err != nil {

		}
		err = light.VerifyAdjacent(i.latestTrustedHeader, untrustedHeader, valSet, 2*168*time.Hour, time.Now(), 10*time.Second)
		if err != nil {
			fmt.Printf("header at height [%d] verify failed\n", height+1)
			i.smartBCHAddrList = rotateList(i.smartBCHAddrList)
			continue
		}
		i.latestTrustedHeader = untrustedHeader
		blk := getBlock(i.smartBCHAddrList, uint64(height+1))
		if bH, tH := blk.Hash(), lb.Hash(); !bytes.Equal(bH, tH) {
			fmt.Printf("block hash not equal signedHeader hash at height [%d]\n", height+1)
			i.smartBCHAddrList = rotateList(i.smartBCHAddrList)
			continue
		}
		//todo: store blk.Txs into db
	}
}

func rotateList(smartbchAddrs []string) []string {
	var newList = []string{smartbchAddrs[len(smartbchAddrs)-1]}
	newList = append(newList, smartbchAddrs[:len(smartbchAddrs)-1]...)
	return newList
}

func (i *Indexer) initLatestTrustedHeader() {
	err := i.db.View(func(txn *badger.Txn) error {
		iterm, err := txn.Get(latestTrustedHeaderKey)
		if err != nil {
			return err
		}
		return iterm.Value(func(val []byte) error {
			var h tmtypes.SignedHeader
			err := tmjson.Unmarshal(val, &h)
			if err != nil {
				return err
			}
			i.latestTrustedHeader = &h
			return nil
		})
	})
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

func (i *Indexer) AddJob(address string, height int64) error {
	if _, exist := i.jobList[address]; exist {
		return jobExistErr
	}
	if len(i.jobList) > maxJobNums {
		return tooManyJobsErr
	}
	j := &Job{indexer: i}
	i.jobList[address] = j
	j.Start()
	return nil
}

func main() {
	indexer := &Indexer{}
	listenURL := flag.String("l", "localhost:8081", " listen address")
	smartBCHAddrListArg := flag.String("b", "http://52.77.220.215:8545,http://18.141.181.5:8545,http://13.250.12.255:8545", "smartbch address list, seperated by comma")
	flag.Parse()
	smartBCHAddrList := strings.Split(*smartBCHAddrListArg, ",")
	if len(smartBCHAddrList) <= 1 {
		log.Fatal("smartbch addresses should at least has two")
	}
	indexer.smartBCHAddrList = smartBCHAddrList
	indexer.listenUrl = *listenURL
	db, err := badger.Open(badger.DefaultOptions("./db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	indexer.db = db
	indexer.initLatestTrustedHeader()
	indexer.chainId = ""
	go indexer.start()

	addHttpHandler(indexer)
	server := http.Server{Addr: *listenURL, ReadTimeout: 3 * time.Second, WriteTimeout: 5 * time.Second}
	fmt.Println("listening ...")
	_ = server.ListenAndServe()
}

func addHttpHandler(indexer *Indexer) {
	http.HandleFunc("/contract", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		addressList := r.URL.Query()["address"]
		if len(addressList) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("address param missed"))
			return
		}
		address := addressList[0]
		bz, err := hex.DecodeString(address)
		if err != nil || len(bz) != 20 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid address param, it format should be hex string without 0x"))
			return
		}
		heights := r.URL.Query()["height"]
		if len(heights) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("height param missed"))
			return
		}
		height := heights[0]
		h, err := strconv.Atoi(height)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid height param format"))
			return
		}
		err = indexer.AddJob(address, int64(h))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
