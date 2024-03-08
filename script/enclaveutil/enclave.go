//go:build linux

package enclaveutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/edgelesssys/ego/enclave"
)

func VerifyEnclaveReportBz(pubKey, reportBz, signerIDBz, uniqueIDBz []byte) error {
	report, err := enclave.VerifyRemoteReport(reportBz)
	if err != nil {
		return err
	}

	if !bytes.Equal(report.SignerID, signerIDBz) {
		return fmt.Errorf("signer-id not match! expected: %x, got: %x", signerIDBz, report.SignerID)
	}
	if !bytes.Equal(report.UniqueID, uniqueIDBz) {
		return fmt.Errorf("unique-id not match! expected: %x, got: %x", uniqueIDBz, report.UniqueID)
	}

	hash := sha256.Sum256(pubKey)
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		return errors.New("report data does not match the pubKey's hash")
	}
	if report.SecurityVersion < 2 {
		return errors.New("invalid security version")
	}
	if binary.LittleEndian.Uint16(report.ProductID) != 0x001 {
		return errors.New("invalid product ID")
	}
	if report.Debug {
		return errors.New("should not open debug mode")
	}

	return nil
}
