package generate

import (
	"alexandreh2ag/go-task/assets"
	"alexandreh2ag/go-task/context"
	"alexandreh2ag/go-task/types"
	"alexandreh2ag/go-task/version"
	"errors"
	"fmt"
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