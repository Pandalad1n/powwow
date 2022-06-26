package test

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
	"time"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	server := exec.CommandContext(
		ctx,
		"go", "run", "../server/.",
		"-listen", "localhost",
		"-port", "8080",
	)
	err := server.Start()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	wisdom := bytes.NewBuffer(nil)
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		client := exec.CommandContext(
			ctx,
			"go", "run", "../client/.",
			"-host", "localhost",
			"-port", "8080",
		)
		client.Stdout = wisdom
		err = client.Start()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		err = client.Wait()
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if wisdom.String() == "" {
		t.Error("Wisdom not returned")
		t.FailNow()
	}
	cancel()
	_ = server.Wait()
}
