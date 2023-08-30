package schedule

import (
	"alexandreh2ag/go-task/context"
	"alexandreh2ag/go-task/types"
	"fmt"
	"github.com/adhocore/gronx"
	"github.com/spf13/afero"
	"os"
	"path"
	"sync"
	"time"
)

const (
	BlocSeparator = "===================="
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

func Run(ctx *context.Context, ref time.Time, noResultPrint bool, resultPath string) []*types.TaskResult {
	var wg sync.WaitGroup
	var mu sync.Mutex

	gron := gronx.New()
	results := []*types.TaskResult{}

	for _, task := range ctx.Config.Scheduled {
		mustRun, err := gron.IsDue(task.CronExpr, ref)
		if err != nil {
			task.Logger.Error(fmt.Sprintf("Scheduled task %s fail to check if must run", task.Id))
		}

		if mustRun {
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
