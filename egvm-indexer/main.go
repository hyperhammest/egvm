package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	jobExistErr    = errors.New("job already exist")
	tooManyJobsErr = errors.New("job nums touch upper limit")
)

const (
	maxJobNums = 1000
)

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
	indexer.chainId = "0x2710"
	indexer.signer = gethtypes.NewEIP155Signer(big.NewInt(10000))
	go indexer.start()

	addHttpHandler(indexer)
	server := http.Server{Addr: *listenURL, ReadTimeout: 3 * time.Second, WriteTimeout: 5 * time.Second}
	fmt.Println("listening ...")
	_ = server.ListenAndServe()
}

func addHttpHandler(indexer *Indexer) {
	// accept new contract scan job
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
