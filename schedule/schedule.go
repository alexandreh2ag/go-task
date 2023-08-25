package schedule

import (
	"alexandreh2ag/go-task/context"
	"alexandreh2ag/go-task/log"
	"alexandreh2ag/go-task/types"
	"fmt"
	"github.com/adhocore/gronx"
	"sync"
	"time"
)

func GetCurrentTime(now time.Time, timezone string) (time.Time, error) {

	if timezone != "" {
		tz, err := time.LoadLocation(timezone)
		if err != nil {
			return now, err
		}
		now = now.In(tz)
	}
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()), nil
}

func Run(ctx *context.Context, ref time.Time) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	gron := gronx.New()
	results := []*types.TaskResult{}

	for _, task := range ctx.Config.Scheduled {
		mustRun, err := gron.IsDue(task.CronExpr, ref)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("Scheduled task %s fail to check if must run", task.Id), log.TaskKey, task.Id)
		}

		if mustRun {
			ctx.Logger.Info(fmt.Sprintf("Scheduled task %s will run", task.Id), log.TaskKey, task.Id)
			wg.Add(1)
			go func(task *types.ScheduledTask) {
				defer wg.Done()

				result := task.Execute()
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}(task)
		} else {
			ctx.Logger.Debug(fmt.Sprintf("Scheduled task %s must not run", task.Id), log.TaskKey, task.Id)
		}
	}
	wg.Wait()
	for _, result := range results {
		ctx.Logger.Debug(fmt.Sprintf("%d: %s", result.Status, result.Output.String()))
	}
}
