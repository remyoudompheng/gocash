package types

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

var guidSeed [16]byte
var guidSeq uint64

func init() {
	_, err := io.ReadFull(rand.Reader, guidSeed[:])
	if err != nil {
		err = fmt.Errorf("failed to initialize GUID generator: %s", err)
		panic(err)
	}
}

type GUID string

// NewGUID returns a pseudo-unique identifier.
func NewGUID() GUID {
	// Gnucash uses a very complicated initialization sequence
	// which we do not reproduce here.
	h := md5.New()
	var buffer [32]byte
	copy(buffer[:16], guidSeed[:])
	binary.LittleEndian.PutUint64(buffer[16:24], uint64(time.Now().UnixNano()))
	binary.LittleEndian.PutUint64(buffer[24:32], atomic.AddUint64(&guidSeq, 1))
	h.Write(buffer[:])
	hash := h.Sum(buffer[:0])
	return GUID(hex.EncodeToString(hash))
}

func (g GUID) Bytes() (b [16]byte, err error) {
	x, err := hex.DecodeString(string(g))
	if len(x) != 16 && err == nil {
		err = fmt.Errorf("invalid GUID of length %s", len(x))
	}
	copy(b[:], x)
	return
}
