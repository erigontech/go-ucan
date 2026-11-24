package envelope

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ucan-wg/go-varsig"
)

type registeredVarsig struct{ payEnc varsig.PayloadEncoding }

const registeredAlgorithm = varsig.Algorithm(0x3a01)

func (v registeredVarsig) Version() varsig.Version                 { return varsig.Version1 }
func (v registeredVarsig) Algorithm() varsig.Algorithm             { return registeredAlgorithm }
func (v registeredVarsig) Hash() varsig.Hash                       { return varsig.HashUnspecified }
func (v registeredVarsig) PayloadEncoding() varsig.PayloadEncoding { return v.payEnc }
func (v registeredVarsig) Encode() []byte                          { return nil }

func TestRegisteredVarsigDecodes(t *testing.T) {
	header := binary.AppendUvarint(nil, varsig.Prefix)
	header = binary.AppendUvarint(header, uint64(varsig.Version1))
	header = binary.AppendUvarint(header, uint64(registeredAlgorithm))
	header = append(header, varsig.EncodePayloadEncoding(varsig.PayloadEncodingDAGCBOR)...)

	_, err := decodeVarsig(header)
	require.ErrorIs(t, err, varsig.ErrUnknownAlgorithm)

	RegisterVarsig(registeredAlgorithm, func(r io.ByteReader) (varsig.Varsig, error) {
		enc, err := varsig.DecodePayloadEncoding(r)
		return registeredVarsig{payEnc: enc}, err
	})

	vsig, err := decodeVarsig(header)
	require.NoError(t, err)
	require.Equal(t, registeredAlgorithm, vsig.Algorithm())
	require.Equal(t, varsig.PayloadEncodingDAGCBOR, vsig.PayloadEncoding())

	_, err = decodeVarsig(varsig.Ed25519(varsig.PayloadEncodingDAGCBOR).Encode())
	require.NoError(t, err)
}
