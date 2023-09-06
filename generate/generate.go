package generate

import (
	"errors"
	"fmt"
	"github.com/alexandreh2ag/go-task/assets"
	"github.com/alexandreh2ag/go-task/context"
	"github.com/alexandreh2ag/go-task/types"
	"github.com/alexandreh2ag/go-task/version"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
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
		err = deleteFile(ctx, outputPath)
		if err != nil {
			return errors.New(fmt.Sprintf("Error when deleting output file: %s", err.Error()))
		}
		return err
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
	}

	tmpl, err := template.New("supervisor.tmpl").Funcs(extraVars).Parse(string(supervisorTemplateContent))
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, ctx.Config.Workers)
}

func generateProgramList(workers types.WorkerTasks) string {
	programs := []string{}
	for _, task := range workers {
		programs = append(programs, task.Id)
	}
	return strings.Join(programs, ",")
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
