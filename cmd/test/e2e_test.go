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
		"-listen", "0.0.0.0",
		"-port", "8888",
	)
	err := server.Start()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	client := exec.CommandContext(
		ctx,
		"go", "run", "../client/.",
		"-port", "8888",
	)
	wisdom := bytes.NewBuffer(nil)
	client.Stdout = wisdom
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		err = client.Start()
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = client.Wait()
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
