package rpc

import (
	"encoding/binary"
	"hash/adler32"
	"io"

	"github.com/pkg/errors"
)

var RPC_MAGIC = [4]byte{'p', 'y', 'x', 'i'}

type Packet struct {
	TotalSize uint32
	Magic     [4]byte
	Payload   []byte
	Checksum  uint32
}

func EncodePacket(w io.Writer, payload []byte) error {
	// len(Magic) + len(Checksum) == 8
	totalsize := uint32(len(RPC_MAGIC) + len(payload) + 4)
	// write totalsize
	binary.Write(w, binary.BigEndian, totalsize)

	sum := adler32.New()
	ww := io.MultiWriter(sum, w)
	// write magic bytes
	binary.Write(ww, binary.BigEndian, RPC_MAGIC)

	// write payload
	ww.Write(payload)

	// calc checksum
	checksum := sum.Sum32()

	// write checksum
	return binary.Write(w, binary.BigEndian, checksum)
}

func DecodePacket(r io.Reader) ([]byte, error) {
	var totalsize uint32
	err := binary.Read(r, binary.BigEndian, &totalsize)
	if err != nil {
		return nil, errors.Wrap(err, "read totalsize")
	}

	if totalsize < 8 {
		return nil, errors.New("totalsize too small")
	}

	sum := adler32.New()
	rr := io.TeeReader(r, sum)

	var magic [4]byte
	binary.Read(rr, binary.BigEndian, &magic)
	if magic != RPC_MAGIC {
		return nil, errors.New("not tiny rpc byte")

	}
	payload := make([]byte, totalsize-8)
	_, err = io.ReadFull(rr, payload)
	if err != nil {
		return nil, errors.Wrap(err, "io readfull payload")
	}
	var checksum uint32
	// read from r, not rr!
	err = binary.Read(r, binary.BigEndian, &checksum)
	if err != nil {
		return nil, errors.Wrap(err, "read checksum")
	}
	if checksum != sum.Sum32() {
		return nil, errors.New("checksum mismatch")
	}
	return payload, nil
}
