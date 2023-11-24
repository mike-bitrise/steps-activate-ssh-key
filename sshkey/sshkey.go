package sshkey

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Agent ...
type Agent interface {
	Start() (string, error)
	Kill() (int, error)
	ListKeys() (int, error)
	AddKey(sshKeyPth, socket string) error
	DeleteKeys() error
}

type defaultAgent struct {
	fileWriter      fileutil.FileWriter
	tempDirProvider pathutil.TempDirProvider
	logger          log.Logger
	cmdFactory      command.Factory
}

// NewAgent ...
func NewAgent(fileWriter fileutil.FileWriter, tempDirProvider pathutil.TempDirProvider, logger log.Logger, cmdFactory command.Factory) Agent {
	return defaultAgent{fileWriter: fileWriter, tempDirProvider: tempDirProvider, logger: logger, cmdFactory: cmdFactory}
}

// Start ...
func (a defaultAgent) Start() (string, error) {
	cmd := a.cmdFactory.Create("ssh-agent", nil, nil)

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnTrimmedOutput()
}

// Kill ...
func (a defaultAgent) Kill() (int, error) {
	// try to kill the agent
	cmd := a.cmdFactory.Create("ssh-agent", []string{"-k"}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnExitCode()
}

// ListKeys ...
func (a defaultAgent) ListKeys() (int, error) {
	cmd := a.cmdFactory.Create("ssh-add", []string{"-l"}, &command.Opts{
		Stderr: os.Stderr,
	})
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnExitCode()
}

// AddKey ...
func (a defaultAgent) AddKey(sshKeyPth, socket string) error {
	var envs []string
	if socket != "" {
		envs = append(envs, "SSH_AUTH_SOCK="+socket)
	}

	cmd := a.cmdFactory.Create("ssh-add", []string{sshKeyPth}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    append(envs, "SSH_ASKPASS=cat"),
		Stdin:  strings.NewReader("nopass"),
	})

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	exitCode, err := cmd.RunAndReturnExitCode()

	if err != nil {
		a.logger.Debugf("Exit code: %s", err)
	}

	if exitCode != 0 {
		a.logger.Errorf("\nExit code: %d", exitCode)
		return fmt.Errorf("failed to add the SSH key to ssh-agent with an empty passphrase")
	}

	return nil
}

// DeleteKeys ...
func (a defaultAgent) DeleteKeys() error {
	// remove all keys from the current agent
	cmd := a.cmdFactory.Create("ssh-add", []string{"-D"}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	a.logger.Println()
	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.Run()
}
