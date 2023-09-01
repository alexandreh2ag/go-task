package schedule

import (
	"alexandreh2ag/go-task/context"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func TestGetScheduleStartCmd_SuccessWithEmptyScheduledTasks(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	fakeClock := clockwork.NewFakeClockAt(time.Date(1970, time.January, 1, 0, 0, 59, 0, time.UTC))
	ctx.Clock = fakeClock
	cmd := GetScheduleStartCmd(ctx)
	go func() {
		err := cmd.Execute()
		assert.NoError(t, err)
	}()

	fakeClock.BlockUntil(1)
	fakeClock.Advance(1 * time.Second)
	time.Sleep(50 * time.Millisecond)

}

func TestGetScheduleStartCmd_ErrorWithWrongTickDuration(t *testing.T) {
	ctx := context.TestContext(io.Discard)
	ctx.Clock = clockwork.NewFakeClockAt(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC))
	cmd := GetScheduleStartCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--" + Tick, "m"})

	err := cmd.Execute()
	assert.Error(t, err)
}
