package cli

import (
	"github.com/alexandreh2ag/go-task/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetScheduleCmd(t *testing.T) {
	ctx := context.TestContext(nil)
	cmd := GetScheduleCmd(ctx)

	assert.Equal(t, 2, len(cmd.Commands()))
}
