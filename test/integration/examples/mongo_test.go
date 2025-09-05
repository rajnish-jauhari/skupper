package examples

import (
	"context"
	"runtime"
	"testing"

	"github.com/skupperproject/skupper/test/integration/examples/mongodb"
)

func TestMongo(t *testing.T) {
	if runtime.GOARCH == "s390x" {
		t.Skip("Skipping TestMongo on s390x architecture")
	}

	mongodb.Run(context.Background(), t, testRunner)
}
