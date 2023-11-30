package schedule

import (
	"fmt"
	"github.com/adhocore/gronx"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/spf13/afero"
	"os"
	"os/signal"
	"path"
	"slices"
	"sync"
	"syscall"
	"time"
)

const (
	BlocSeparator = "===================="
)

var (
	cronExprNextTick = "*/%d * * * *"
	tickUnit         = time.Minute
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

func Start(ctx *context.Context, tick int, timezone string, taskFilter []string, noResultPrint bool, resultPath string) error {
	var refTime time.Time
	firstRun := true

	now := ctx.Clock.Now()
	nextTick, err := gronx.NextTickAfter(fmt.Sprintf(cronExprNextTick, tick), now, false)
	if err != nil {
		return fmt.Errorf("could not calculate next tick of expr %s: %v", fmt.Sprintf(cronExprNextTick, tick), err)
	}
	ctx.Logger.Info(fmt.Sprintf("next tick: %v", nextTick.Format("2006-01-02T15:04:05")))
	ctx.Clock.Sleep(nextTick.Sub(now))

	refTime, err = GetCurrentTime(ctx.Clock.Now(), timezone)
	if err != nil {
		return err
	}

	ticker := ctx.Clock.NewTicker(time.Duration(tick) * tickUnit)
	defer ticker.Stop()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		if firstRun {
			firstRun = false
			ctx.Logger.Debug("first tick")
			Run(ctx, refTime, taskFilter, false, noResultPrint, resultPath)
		}
		select {
		case <-ticker.Chan():
			ctx.Logger.Debug("tick")
			// can ignore error because schedule.GetCurrentTime used at top
			refTime, _ = GetCurrentTime(ctx.Clock.Now(), timezone)

			Run(ctx, refTime, taskFilter, false, noResultPrint, resultPath)

		case sig := <-sigs:
			ctx.Logger.Info(fmt.Sprintf("%s signal received, exiting...", sig.String()))
			return nil

		case <-ctx.Done():
			ctx.Logger.Info(fmt.Sprintf("stop asked by app, exiting..."))
			return nil
		}
	}
}

func Run(ctx *context.Context, ref time.Time, taskFilter []string, force bool, noResultPrint bool, resultPath string) []*types.TaskResult {
	var wg sync.WaitGroup
	var mu sync.Mutex

	gron := gronx.New()
	results := []*types.TaskResult{}

	for _, task := range ctx.Config.Scheduled {
		if len(taskFilter) != 0 && !slices.Contains(taskFilter, task.Id) {
			ctx.Logger.Info(fmt.Sprintf("Task %s skipped due to the filter %v", task.Id, taskFilter))
			continue
		}

		mustRun, err := gron.IsDue(task.CronExpr, ref)
		if err != nil {
			task.Logger.Error(fmt.Sprintf("Scheduled task %s fail to check if must run", task.Id))
		}

		if mustRun || force {
			task.Logger.Info(fmt.Sprintf("Scheduled task %s will run", task.Id))
			wg.Add(1)
			go func(task *types.ScheduledTask, noResultPrint bool, resultPath string) {
				defer wg.Done()

				result := task.Execute()
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				output := FormatTaskResult(result)

				if !noResultPrint {
					fmt.Println(output)
				}

				if resultPath != "" {
					err = WriteToLogFile(ctx, resultPath, output)
					if err != nil {
						task.Logger.Error(err.Error())
					}
				}

			}(task, noResultPrint, resultPath)
		} else {
			task.Logger.Debug(fmt.Sprintf("Scheduled task %s must not run", task.Id))
		}
	}
	wg.Wait()

	return results
}

func FormatTaskResult(result *types.TaskResult) string {
	var outputStr = ""
	var errorStr = ""

	if result.Output.String() != "" {
		outputStr = fmt.Sprintf("output:\n%s\n", result.Output.String())
	}

	if result.Error != nil {
		errorStr = fmt.Sprintf("Due to the following error: %s\n", result.Error.Error())
	}
	return fmt.Sprintf(
		"%s\nTask %s finish with status '%s'\nStart at %s, finish at %s (%s)\n%s%s%s\n",
		BlocSeparator,
		result.Task.Id,
		result.StatusString(),
		result.StartAt.Format("2006-01-02T15:04:05"),
		result.FinishAt.Format("2006-01-02T15:04:05"),
		result.FinishAt.Sub(result.StartAt),
		outputStr,
		errorStr,
		BlocSeparator,
	)
}

func WriteToLogFile(ctx *context.Context, resultPath string, data string) error {
	afs := &afero.Afero{Fs: ctx.Fs}
	if ok, _ := afs.DirExists(path.Dir(resultPath)); !ok {
		err := ctx.Fs.MkdirAll(path.Dir(resultPath), 0770)
		if err != nil {
			return fmt.Errorf("failed to create log path %s with error %s", resultPath, err.Error())
		}
	}

	file, err := ctx.Fs.OpenFile(resultPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o660)
	if err != nil {
		return fmt.Errorf("failed to open log path %s with error %s", resultPath, err.Error())
	}

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write in log path %s with error %s", resultPath, err.Error())
	}
	return nil
}
