package testdata

import (
	"testing"
	"time"
)

func TestOne(t *testing.T) {
	time.Sleep(4 * time.Second)
}

func TestTwo(t *testing.T) {
	t.Log("this is some output\nline 2")
	time.Sleep(3 * time.Second)
	t.Log("end of test")
}

func TestThree(t *testing.T) {
	time.Sleep(1 * time.Second)
	t.Skipf("skipping this test...")
}

func TestFour(t *testing.T) {
	time.Sleep(4 * time.Second)
	t.Log("this test failed for some reason...")
	t.FailNow()
}
