package hashcash

import (
	"fmt"
	"testing"
)

func TestChallenge(t *testing.T) {
	c := NewChallenge(18)
	solution := c.Solve()
	if !c.Verify(solution) {
		t.Fail()
	}
}

func TestChallengeFalse(t *testing.T) {
	c := NewChallenge(18)
	solution := []byte{12}
	if c.Verify(solution) {
		t.Fail()
	}
}

func BenchmarkChallenge(b *testing.B) {
	difficulties := [...]uint32{1, 5, 10, 15, 20}
	for _, d := range difficulties {
		b.Run(fmt.Sprintf("Testing %v difficulty", d), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := NewChallenge(d)
				solution := c.Solve()
				if !c.Verify(solution) {
					b.Fail()
				}
			}
		})
	}
}
