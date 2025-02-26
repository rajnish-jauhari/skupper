//go:build integration || examples
// +build integration examples

package examples

import (
	"context"
	"testing"
	"runtime"

	"github.com/skupperproject/skupper/test/integration/examples/bookinfo"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestBookinfo(t *testing.T) {
	// Check if the architecture is s390x
	if runtime.GOARCH == "s390x" {
		t.Skip("Skipping test on s390x architecture")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bookinfo.Run(ctx, t, testRunner)
}
