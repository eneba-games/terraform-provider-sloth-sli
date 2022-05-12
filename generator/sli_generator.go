package generator

import (
	"bytes"
	"fmt"
	"os/exec"
)

type SLI struct {
	slothExecPath string
}

func NewSLI(slothExecPath string) *SLI {
	return &SLI{
		slothExecPath: slothExecPath,
	}
}

func (g *SLI) Generate(inputFile string) (string, error) {
	if err := g.validate(inputFile); err != nil {
		return "", fmt.Errorf("")
	}

	out, err := g.generate(inputFile)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (g *SLI) validate(inputFile string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(g.slothExecPath, "validate", "-i", inputFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return formatCmdErr("SLI config validation failed", stderr, err)
	}

	return nil
}

func (g *SLI) generate(inputFile string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(g.slothExecPath, "generate", "-i", inputFile)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", formatCmdErr("failed to generate SLI recording & alerting rules", stderr, err)
	}

	return stdout.String(), nil
}

func formatCmdErr(msg string, stderr bytes.Buffer, err error) error {
	if stderr.Len() > 0 {
		return fmt.Errorf("%s: %s", msg, stderr.String())
	}

	return fmt.Errorf("%s: %w", msg, err)
}
