package derive

import (
<<<<<<< HEAD
	"bytes"
	"testing"

	"github.com/ethereum-optimism/optimism/op-node/rollup"

=======
	"testing"

>>>>>>> v0.5.23
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestBatchRoundTrip(t *testing.T) {
	batches := []*BatchData{
		{
			BatchV1: BatchV1{
<<<<<<< HEAD
				Epoch:        0,
=======
				EpochNum:     0,
>>>>>>> v0.5.23
				Timestamp:    0,
				Transactions: []hexutil.Bytes{},
			},
		},
		{
			BatchV1: BatchV1{
<<<<<<< HEAD
				Epoch:        1,
=======
				EpochNum:     1,
>>>>>>> v0.5.23
				Timestamp:    1647026951,
				Transactions: []hexutil.Bytes{[]byte{0, 0, 0}, []byte{0x76, 0xfd, 0x7c}},
			},
		},
	}

	for i, batch := range batches {
		enc, err := batch.MarshalBinary()
		assert.NoError(t, err)
		var dec BatchData
		err = dec.UnmarshalBinary(enc)
		assert.NoError(t, err)
		assert.Equal(t, batch, &dec, "Batch not equal test case %v", i)
	}
<<<<<<< HEAD
	var buf bytes.Buffer
	err := EncodeBatches(&rollup.Config{}, batches, &buf)
	assert.NoError(t, err)
	out, err := DecodeBatches(&rollup.Config{}, &buf)
	assert.NoError(t, err)
	assert.Equal(t, batches, out)
=======
>>>>>>> v0.5.23
}
