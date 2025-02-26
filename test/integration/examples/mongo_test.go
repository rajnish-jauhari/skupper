//go:build integration || examples
// +build integration examples

package examples

import (
	"context"
	"runtime"
	"testing"

	"github.com/skupperproject/skupper/test/integration/examples/mongodb"
)

func TestMongo(t *testing.T) {
	// Check if the architecture is s390x
	if runtime.GOARCH == "s390x" {
		t.Skip("Skipping test on s390x architecture")
	}

	mongodb.Run(context.Background(), t, testRunner)
}