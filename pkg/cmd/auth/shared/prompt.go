package shared

import (
	"fmt"
	"os"
)

// Prompt represents a command-line prompt
type Prompt struct {
	IOStreams *iostreams.IOStreams
}

// NewPrompt creates a new Prompt instance
func NewPrompt(io *iostreams.IOStreams) *Prompt {
	return &Prompt{
		IOStreams: io,
	}
}

// Input prompts the user for input
func (p *Prompt) Input(prompt string) (string, error) {
	fmt.Fprint(p.IOStreams.Out, prompt)
	
	var input string
	_, err := fmt.Fscanln(p.IOStreams.In, &input)
	if err != nil {
		return "", err
	}
	
	return input, nil
}

// Confirm prompts the user for a yes/no confirmation
func (p *Prompt) Confirm(prompt string, def bool) (bool, error) {
	choices := "y/N"
	if def {
		choices = "Y/n"
	}
	
	fmt.Fprintf(p.IOStreams.Out, "%s [%s] ", prompt, choices)
	
	var response string
	_, err := fmt.Fscanln(p.IOStreams.In, &response)
	if err != nil {
		return def, nil
	}
	
	switch response {
	case "y", "Y", "yes", "YES", "Yes":
		return true, nil
	case "n", "N", "no", "NO", "No":
		return false, nil
	default:
		return def, nil
	}
}

// Password prompts the user for a password (no echo)
func (p *Prompt) Password(prompt string) (string, error) {
	fmt.Fprint(p.IOStreams.Out, prompt)
	
	// In a real implementation, we would use a library like "golang.org/x/term"
	// For demo purposes, we'll just read normally
	var password string
	_, err := fmt.Fscanln(p.IOStreams.In, &password)
	if err != nil {
		return "", err
	}
	
	return password, nil
}

// Select prompts the user to choose from a list
func (p *Prompt) Select(prompt string, options []string, def int) (int, error) {
	for i, option := range options {
		defaultStr := ""
		if i == def {
			defaultStr = " (default)"
		}
		fmt.Fprintf(p.IOStreams.Out, "%d) %s%s\n", i+1, option, defaultStr)
	}
	
	fmt.Fprintf(p.IOStreams.Out, "%s [%d]: ", prompt, def+1)
	
	var choice int
	_, err := fmt.Fscanln(p.IOStreams.In, &choice)
	if err != nil {
		return def, nil
	}
	
	if choice < 1 || choice > len(options) {
		return def, nil
	}
	
	return choice - 1, nil
}

// MultiSelect prompts the user to select multiple options
func (p *Prompt) MultiSelect(prompt string, options []string, defs []bool) ([]bool, error) {
	result := make([]bool, len(options))
	
	for i, option := range options {
		defStr := ""
		if i < len(defs) && defs[i] {
			defStr = " (selected)"
			result[i] = true
		}
		
		fmt.Fprintf(p.IOStreams.Out, "%d) %s%s\n", i+1, option, defStr)
	}
	
	fmt.Fprintf(p.IOStreams.Out, "%s (comma-separated numbers, or 'all'): ", prompt)
	
	var input string
	_, err := fmt.Fscanln(p.IOStreams.In, &input)
	if err != nil {
		return result, nil
	}
	
	// For demo, just return defaults
	return result, nil
}

// Checkmark prints a checkmark for success
func (p *Prompt) Checkmark() {
	fmt.Fprintln(p.IOStreams.Out, "✓")
}

// Error prints an error message
func (p *Prompt) Error(message string) {
	fmt.Fprintf(p.IOStreams.ErrOut, "Error: %s\n", message)
}

// Info prints an info message
func (p *Prompt) Info(message string) {
	fmt.Fprintf(p.IOStreams.Out, "Info: %s\n", message)
}

// Warn prints a warning message
func (p *Prompt) Warn(message string) {
	fmt.Fprintf(p.IOStreams.ErrOut, "Warning: %s\n", message)
}

// Debug prints a debug message
func (p *Prompt) Debug(message string) {
	if os.Getenv("DEBUG") != "" {
		fmt.Fprintf(p.IOStreams.ErrOut, "Debug: %s\n", message)
	}
}