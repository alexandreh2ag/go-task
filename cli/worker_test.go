package cli

import (
	"alexandreh2ag/go-task/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWorkerCmd(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetWorkerCmd(ctx)

	assert.Equal(t, 1, len(cmd.Commands()))
}
