package types

import (
	"fmt"
	"io"

	"github.com/bytom-gm/crypto/sm3"
	"github.com/bytom-gm/encoding/blockchain"
	"github.com/bytom-gm/errors"
	"github.com/bytom-gm/protocol/bc"
)

// OutputCommitment contains the commitment data for a transaction output.
type OutputCommitment struct {
	bc.AssetAmount
	VMVersion      uint64
	ControlProgram []byte
}

func (oc *OutputCommitment) writeExtensibleString(w io.Writer, suffix []byte, assetVersion uint64) error {
	_, err := blockchain.WriteExtensibleString(w, suffix, func(w io.Writer) error {
		return oc.writeContents(w, suffix, assetVersion)
	})
	return err
}

func (oc *OutputCommitment) writeContents(w io.Writer, suffix []byte, assetVersion uint64) (err error) {
	if assetVersion == 1 {
		if _, err = oc.AssetAmount.WriteTo(w); err != nil {
			return errors.Wrap(err, "writing asset amount")
		}
		if _, err = blockchain.WriteVarint63(w, oc.VMVersion); err != nil {
			return errors.Wrap(err, "writing vm version")
		}
		if _, err = blockchain.WriteVarstr31(w, oc.ControlProgram); err != nil {
			return errors.Wrap(err, "writing control program")
		}
	}
	if len(suffix) > 0 {
		_, err = w.Write(suffix)
	}
	return errors.Wrap(err, "writing suffix")
}

func (oc *OutputCommitment) readFrom(r *blockchain.Reader, assetVersion uint64) (suffix []byte, err error) {
	return blockchain.ReadExtensibleString(r, func(r *blockchain.Reader) error {
		if assetVersion == 1 {
			if err := oc.AssetAmount.ReadFrom(r); err != nil {
				return errors.Wrap(err, "reading asset+amount")
			}
			oc.VMVersion, err = blockchain.ReadVarint63(r)
			if err != nil {
				return errors.Wrap(err, "reading VM version")
			}
			if oc.VMVersion != 1 {
				return fmt.Errorf("unrecognized VM version %d for asset version 1", oc.VMVersion)
			}
			oc.ControlProgram, err = blockchain.ReadVarstr31(r)
			return errors.Wrap(err, "reading control program")
		}
		return nil
	})
}

// Hash convert suffix && assetVersion to bc.Hash
func (oc *OutputCommitment) Hash(suffix []byte, assetVersion uint64) (outputhash bc.Hash) {
	h := sm3.New()
	oc.writeExtensibleString(h, suffix, assetVersion)
	var b32 [32]byte
	copy(b32[:], h.Sum(nil))
	outputhash = bc.NewHash(b32)
	return outputhash
}
