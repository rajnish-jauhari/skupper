//go:build integration || examples
// +build integration examples

package examples

import (
	"context"
	"testing"
	"runtime"

	"github.com/skupperproject/skupper/test/integration/examples/mongodb"
)

func TestMongo(t *testing.T) {
	if runtime.GOARCH == "s390x" {
		t.Skip("Skipping test on s390x architecture")
	}
	mongodb.Run(context.Background(), t, testRunner)
}
