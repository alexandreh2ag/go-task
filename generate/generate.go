package generate

import (
	"errors"
	"fmt"
	"github.com/alexandreh2ag/go-task/assets"
	"github.com/alexandreh2ag/go-task/condition"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/env"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/alexandreh2ag/go-task/version"
	"github.com/spf13/afero"
	"golang.org/x/exp/maps"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	FormatSupervisor = "supervisor"
)

func Generate(ctx *context.Context, outputPath string, format string, groupName string) error {
	err := checkDir(ctx, outputPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error with outputh dir: %s", err.Error()))
	}

	if len(ctx.Config.Workers) == 0 {
		afs := &afero.Afero{Fs: ctx.Fs}
		if ok, _ := afs.Exists(outputPath); ok {
			err = deleteFile(ctx, outputPath)
			if err != nil {
				return errors.New(fmt.Sprintf("Error when deleting output file: %s", err.Error()))
			}
		}
		return nil
	}

	outputFile, err := ctx.Fs.Create(outputPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error with output file: %s", err.Error()))
	}

	switch format {
	case FormatSupervisor:
		return templateSupervisorFile(ctx, outputFile, groupName)
	default:
		return errors.New(fmt.Sprintf("Error with unsupported format %s", format))
	}
}

// check if directory indicated in path exist
func checkDir(ctx *context.Context, path string) error {
	dir := filepath.Dir(path)
	_, err := ctx.Fs.Stat(dir)
	if err != nil {
		return err
	}
	return nil
}

func templateSupervisorFile(ctx *context.Context, writer io.Writer, groupName string) error {
	supervisorTemplateContent, err := fs.ReadFile(assets.TemplateFiles, "templates/supervisor.tmpl")
	if err != nil {
		return errors.New(fmt.Sprintf("Error with template file: %s", err.Error()))
	}

	extraVars := template.FuncMap{
		"now":       time.Now,
		"version":   version.GetFormattedVersion,
		"groupName": func() string { return groupName },
		"programs":  generateProgramList,
		"envs":      generateEnvVars,
	}

	tmpl, err := template.New("supervisor.tmpl").Funcs(extraVars).Parse(string(supervisorTemplateContent))
	if err != nil {
		return err
	}
	templatedWorkers := types.WorkerTasks{}

	for _, worker := range ctx.Config.Workers {
		result, err := condition.EvalExpression(worker.Expression, worker.Envs)
		if err != nil {
			return fmt.Errorf("can't evaluate expression for task '%s': %v", worker.Id, err)
		}
		if result {
			templatedWorkers = append(templatedWorkers, worker)
		} else {
			ctx.Logger.Info(fmt.Sprintf("skipping task '%s': expression false", worker.Id))
		}
	}

	return tmpl.Execute(writer, templatedWorkers)
}

func generateProgramList(workers types.WorkerTasks) string {
	programs := []string{}
	for _, task := range workers {
		programs = append(programs, task.PrefixedName())
	}
	return strings.Join(programs, ",")
}

func generateEnvVars(worker types.WorkerTask) string {
	envVars := []string{}

	// ordering key to have deterministic results
	keys := maps.Keys(worker.Envs)
	sort.Strings(keys)

	for _, varName := range keys {
		envVars = append(envVars, fmt.Sprintf(`%s="%s"`, varName, os.Expand(worker.Envs[varName], env.GetEnvVars(worker.Envs))))
	}
	return strings.Join(envVars, ",")
}

func deleteFile(ctx *context.Context, path string) error {
	err := ctx.Fs.Remove(path)
	if err != nil {
		if strings.Contains(err.Error(), "file does not exist") {
			return nil
		}
	}
	return err
}
