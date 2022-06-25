package hashcash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
)

type Challenge struct {
	Digest     []byte
	Difficulty uint32
}

func NewChallenge(difficulty uint32) Challenge {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return Challenge{Digest: b, Difficulty: difficulty}
}

func (c Challenge) Solve() []byte {
	var i uint64
	buf := make([]byte, 8)
	for {
		binary.LittleEndian.PutUint64(buf, i)
		res := c.Verify(buf)
		if res {
			return buf
		}
		i++
	}
}

func (c Challenge) Verify(solution []byte) bool {
	hasher := sha256.New()
	_, err := io.Copy(hasher, bytes.NewReader(c.Digest))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(hasher, bytes.NewReader(solution))
	if err != nil {
		panic(err)
	}

	hash := hasher.Sum(nil)
	return verifyZeros(hash, c.Difficulty)
}

func verifyZeros(data []byte, zeros uint32) bool {
	getBit := func(from byte, pos uint8) bool {
		return from&(1<<(7-pos)) != 0
	}
	var count uint32
	for _, bt := range data {
		btPos := uint8(0)
		for btPos <= 7 {
			if getBit(bt, btPos) {
				return false
			}
			count++
			btPos++
			if count >= zeros {
				return true
			}
		}
	}
	return false
}
