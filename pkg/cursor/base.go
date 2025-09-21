package cursor

import (
	"context"
	"errors"
	"io"
	"os/exec"

	log "git.oceantim.com/backend/packages/golang/go-logger"
)

type OutputFormat string

const OutFormatText OutputFormat = "text"
const OutFormatJSON OutputFormat = "json"

type Model string

const (
	ModelChatGPT5 Model = "gpt-5"
	ModelOpus41   Model = "opus-4.1"
	ModelSonnet4        = "sonnet-4"
	ModelGrok           = "grok"
)

const command = "cursor-agent"

type Cursor struct {
	args   []string
	logger log.Logger
	prompt string
}

type ExecConfig struct {
	Dir    string
	Stdout io.Writer
	Stderr io.Writer
}

func New() *Cursor {
	return &Cursor{
		args: []string{
			"-p",
		},
	}
}

func (r *Cursor) Interactive() *Cursor {
	r.args = append(r.args, "--force")

	return r
}

func (r *Cursor) OutputFormat(OutputFormat OutputFormat) *Cursor {
	r.args = append(r.args, "--output-format", string(OutputFormat))

	return r
}

func (r *Cursor) Model(model Model) *Cursor {
	r.args = append(r.args, "--model ", string(model))

	return r
}

func (r *Cursor) Prompt(prompt string) *Cursor {
	r.prompt = prompt

	return r
}

func (r *Cursor) Execute(ctx context.Context, config *ExecConfig) error {
	cmd := exec.CommandContext(ctx, command, append(r.args, r.prompt)...)
	cmd.Dir = config.Dir
	cmd.Stdout = config.Stdout
	cmd.Stderr = config.Stderr

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			r.logger.Error("command not successful", log.J{
				"command": command,
				"args":    r.args,
				"stderr":  exitErr.Stderr,
			})
		}
		return err
	}

	return nil
}
