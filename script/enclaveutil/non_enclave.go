//go:build darwin

package enclaveutil

func VerifyEnclaveReportBz(pubKey, reportBz, signerIDBz, uniqueIDBz []byte) error {
	return nil
}
