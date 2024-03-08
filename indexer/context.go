package indexer

import (
	"bytes"
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"hash"
	"io"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	InputLogFileName  = "in.txt"
	OutputLogFileName = "out.txt"
)

var Signer = gethtypes.NewEIP155Signer(big.NewInt(10000))

type SimpleTx struct {
	Sender string `json:"sender"`
	Data   []byte `json:"data"`
}

type SimpleEvent struct {
	Topics []string `json:"topics"`
	Data   []string `json:"data"`
}

type BlockForContract struct {
	Height    int64         `json:"height"`
	Timestamp int64         `json:"time"`
	TxList    []SimpleTx    `json:"txs"`
	EventList []SimpleEvent `json:"events"`
}

type Context struct {
	VM              *goja.Runtime             `json:"-"`
	HandleBlockFunc goja.Callable             `json:"-"`
	InputLog        *os.File                  `json:"-"`
	OutputLog       *os.File                  `json:"-"`
	InputLogHasher  hash.Hash                 `json:"-"`
	OutputLogHasher hash.Hash                 `json:"-"`
	InputLogHash    [32]byte                  `json:"ihash"`
	OutputLogHash   [32]byte                  `json:"ohash"`
	InputLogSize    int64                     `json:"isize"`
	OutputLogSize   int64                     `json:"osize"`
	UserToNonce     map[common.Address]uint64 `json:"nonces"`
	TrustedHeader   *tmtypes.SignedHeader     `json:"header"`
}

func newContext(vm *goja.Runtime) *Context {
	c := &Context{VM: vm}

	handleBlockFunc, ok := goja.AssertFunction(vm.Get("handleBlock"))
	if !ok {
		log.Fatalf("cannot find the 'handleBlock' function")
	}
	c.HandleBlockFunc = handleBlockFunc

	return c
}

func NewContext(vm *goja.Runtime) *Context {
	c := newContext(vm)

	f, err := os.OpenFile(InputLogFileName, os.O_EXCL, 0644)
	if err != nil {
		log.Fatalf("cannot open %s: %v", InputLogFileName, err)
	}
	c.InputLog = f

	f, err = os.OpenFile(OutputLogFileName, os.O_EXCL, 0644)
	if err != nil {
		log.Fatalf("cannot open %s: %v", OutputLogFileName, err)
	}
	c.OutputLog = f

	return c
}

func NewContextFromRecoverFile(recoveryFileName string, vm *goja.Runtime) *Context {
	c := newContext(vm)

	fileData, err := os.ReadFile(recoveryFileName)
	if err != nil {
		log.Fatalf("cannot open %s: %v", recoveryFileName, err)
	}
	if len(fileData) <= 65 {
		log.Fatalf("incorrect length for recovery file")
	}
	sig := fileData[:65]
	jsonData := fileData[65:]
	hash := sha256.Sum256(jsonData)
	ok := crypto.VerifySignature(KGClient.PubKeyBz, hash[:], sig[:64])
	if !ok {
		log.Fatalf("verify recovery sig failed!")
	}

	err = json.Unmarshal(jsonData, c)
	if err != nil {
		panic(err)
	}
	var inHash, outHash []byte
	inHash, c.InputLogHasher = truncateAndRunHasher(InputLogFileName, c.InputLogSize)
	if !bytes.Equal(inHash, c.InputLogHash[:]) {
		panic("Invalid Hash for "+InputLogFileName)
	}
	outHash, c.OutputLogHasher = truncateAndRunHasher(OutputLogFileName, c.OutputLogSize)
	if !bytes.Equal(outHash, c.OutputLogHash[:]) {
		panic("Invalid Hash for "+OutputLogFileName)
	}

	c.replayInput()
	return c
}

func truncateAndRunHasher(fname string, size int64) ([]byte, hash.Hash) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = f.Truncate(size)
	if err != nil {
		panic(err)
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	return h.Sum(nil), h
}

func (c *Context) replayInput() {
	f, err := os.Open(InputLogFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		//fmt.Println(line)
		c.HandleBlockFunc(goja.Undefined(), c.VM.ToValue(line))
		h.Write([]byte(line))
	}

	inHash := h.Sum(nil)
	if !bytes.Equal(inHash, c.InputLogHash[:]) {
		panic("Invalid Hash for "+InputLogFileName)
	}
}

func (c *Context) DumpRecoveryFile(recoveryFileName string, header *tmtypes.SignedHeader) {
	c.TrustedHeader = header
	copy(c.InputLogHash[:], c.InputLogHasher.Sum(nil))
	copy(c.OutputLogHash[:], c.OutputLogHasher.Sum(nil))
	c.InputLog.Sync()
	c.OutputLog.Sync()
	fileInfo, err := c.InputLog.Stat()
	if err != nil {
		panic(err)
	}
	c.InputLogSize = fileInfo.Size()
	fileInfo, err = c.OutputLog.Stat()
	if err != nil {
		panic(err)
	}
	c.OutputLogSize = fileInfo.Size()
	bz, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(bz)
	sig, err := crypto.Sign(hash[:], KGClient.PrivKey.ToECDSA())
	if err != nil {
		panic(err)
	}
	out := append(sig, bz...)
	err = os.WriteFile(recoveryFileName, out, 0600)
	if err != nil {
		panic(err)
	}
}

func (c *Context) HandleBlock(contract common.Address, blk *tmtypes.Block) {
	b := getBlockForContract(contract, blk)
	line, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	c.InputLogHasher.Write([]byte(line))
	c.InputLogHasher.Write([]byte("\n"))
	c.InputLog.Write([]byte(line))
	c.InputLog.Write([]byte("\n"))
	res, _ := c.HandleBlockFunc(goja.Undefined(), c.VM.ToValue(line))
	var outStr string
	c.VM.ExportTo(res, &outStr)
	c.OutputLogHasher.Write([]byte(outStr))
	c.OutputLogHasher.Write([]byte("\n"))
	c.OutputLog.Write([]byte(outStr))
	c.OutputLog.Write([]byte("\n"))
}

func getBlockForContract(contract common.Address, blk *tmtypes.Block) (b BlockForContract) {
	b.Height = blk.Height
	b.Timestamp = blk.Time.UnixMilli()
	for _, txRawData := range blk.Txs {
		tx := &gethtypes.Transaction{}
		err := tx.DecodeRLP(rlp.NewStream(bytes.NewReader(txRawData), 0))
		if err != nil {
			continue
		}
		if *tx.To() == contract {
			sender, err := Signer.Sender(tx)
			if err != nil {
				continue
			}
			tx := SimpleTx{
				Sender: sender.String(),
				Data:   tx.Data(),
			}
			b.TxList = append(b.TxList, tx)
		}
	}
	return
}

