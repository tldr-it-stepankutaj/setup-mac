package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

// Prompt provides interactive prompts
type Prompt struct {
	Interactive bool
}

// NewPrompt creates a new Prompt
func NewPrompt(interactive bool) *Prompt {
	return &Prompt{
		Interactive: interactive,
	}
}

// Confirm asks for yes/no confirmation.
//
// Implemented with bufio rather than promptui.IsConfirm because promptui ignores
// the Default field in confirm mode and treats Enter as ErrAbort, which silently
// flipped "Proceed? [Y/n]" to "no" whenever the user just hit return.
func (p *Prompt) Confirm(label string, defaultVal bool) (bool, error) {
	if !p.Interactive {
		return defaultVal, nil
	}

	suffix := "[y/N]"
	if defaultVal {
		suffix = "[Y/n]"
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "%s %s: ", label, suffix)

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return defaultVal, nil
			}
			return defaultVal, err
		}

		switch strings.ToLower(strings.TrimSpace(line)) {
		case "":
			return defaultVal, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		fmt.Fprintln(os.Stderr, "Please answer y or n.")
	}
}

// Input asks for text input
func (p *Prompt) Input(label, defaultVal string) (string, error) {
	if !p.Interactive {
		return defaultVal, nil
	}

	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultVal,
	}

	return prompt.Run()
}

// InputRequired asks for required text input
func (p *Prompt) InputRequired(label string) (string, error) {
	if !p.Interactive {
		return "", fmt.Errorf("interactive mode required for this input")
	}

	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("input cannot be empty")
			}
			return nil
		},
	}

	return prompt.Run()
}

// Select asks to select from a list
func (p *Prompt) Select(label string, items []string) (int, string, error) {
	if !p.Interactive {
		if len(items) > 0 {
			return 0, items[0], nil
		}
		return -1, "", fmt.Errorf("no items to select")
	}

	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  10,
	}

	return prompt.Run()
}

// SelectWithDescription asks to select from a list with descriptions
func (p *Prompt) SelectWithDescription(label string, items []SelectItem) (int, *SelectItem, error) {
	if !p.Interactive {
		if len(items) > 0 {
			return 0, &items[0], nil
		}
		return -1, nil, fmt.Errorf("no items to select")
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Description | faint }})",
		Inactive: "  {{ .Name | white }} ({{ .Description | faint }})",
		Selected: "\U00002705 {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return -1, nil, err
	}

	return idx, &items[idx], nil
}

// SelectItem represents an item with description
type SelectItem struct {
	Name        string
	Description string
	Value       string
}

// Password asks for password input (hidden)
func (p *Prompt) Password(label string) (string, error) {
	if !p.Interactive {
		return "", fmt.Errorf("interactive mode required for password input")
	}

	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}

	return prompt.Run()
}
