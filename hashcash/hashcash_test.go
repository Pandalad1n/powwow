package hashcash

import "testing"

func TestChallenge(t *testing.T) {
	c := NewChallenge(18)
	solution := c.Solve()
	if !c.Verify(solution) {
		t.Fail()
	}
}
