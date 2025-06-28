package io

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/mod/modfile"
)

var (
	// ErrGoModNotFound is an error returned when the go.mod file is not found.
	ErrGoModNotFound = errors.New("could not find go.mod env (outside a go module?)")
)

// GetGoModFile reads and retrieves the go.mod file.
func GetGoModFile() (*modfile.File, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	gomodenv := &bytes.Buffer{}
	cmd.Stdout = gomodenv

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("could not get go.mod path: %w", err)
	}

	if gomodenv.String() == "/dev/null" {
		return nil, ErrGoModNotFound
	}

	content, err := os.ReadFile("go.mod")
	if err != nil {
		return nil, fmt.Errorf("could not read go.mod file: %w", err)
	}

	file, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse go.mod file: %w", err)
	}

	return file, nil
}
