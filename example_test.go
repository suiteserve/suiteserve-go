package go_runner

import (
	"fmt"
	"testing"
	"time"
)

func TestExample1(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
	fmt.Println("EXAMPLE 1")
}

func TestExample2(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}

func TestExample3(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}

func TestExample4(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}

func TestExample5(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}

func TestExample6(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}
