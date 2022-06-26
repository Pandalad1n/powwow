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

// Solve can be done in parallel to speed it up.
// But we keep it simple.
func (c Challenge) Solve() []byte {
	var i uint64
	solution := make([]byte, 8)
	for {
		binary.LittleEndian.PutUint64(solution, i)
		if c.Verify(solution) {
			return solution
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

// verifyZeros returns true if `data` starts with a `zeros` number of unset bits.
func verifyZeros(data []byte, zeros uint32) bool {
	getBit := func(from byte, pos uint8) bool {
		// Creating a mask by shifting 1 to required position.
		// Then applying that mask to `from`.
		// If resulting value is not 0 then the corresponding bit is set.
		return from&(1<<(7-pos)) != 0
	}
	var count uint32
	for _, byteVal := range data {
		bitPos := uint8(0)
		for bitPos <= 7 {
			if getBit(byteVal, bitPos) {
				return false
			}
			count++
			bitPos++
			if count >= zeros {
				return true
			}
		}
	}
	return false
}
