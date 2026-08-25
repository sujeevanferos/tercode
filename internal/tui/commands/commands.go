// Package commands implements slash commands (/model, /clear, /diff, /config, /help, etc.).
package commands

import (
	"fmt"
	"strings"
)

// Command represents an executable slash command.
type Command struct {
	Name        string
	Description string
	Execute     func(args []string) string
}

// Registry stores all built-in and user-defined slash commands.
type Registry struct {
	commands map[string]Command
}

// New creates an initialized CommandRegistry.
func New() *Registry {
	r := &Registry{commands: make(map[string]Command)}
	r.registerBuiltins()
	return r
}

func (r *Registry) registerBuiltins() {
	r.Register(Command{
		Name:        "help",
		Description: "Show available slash commands",
		Execute: func(args []string) string {
			var sb strings.Builder
			sb.WriteString("Available Commands:\n")
			for _, cmd := range r.commands {
				sb.WriteString(fmt.Sprintf("  /%-12s %s\n", cmd.Name, cmd.Description))
			}
			return sb.String()
		},
	})
	r.Register(Command{
		Name:        "clear",
		Description: "Clear conversation screen",
		Execute: func(args []string) string {
			return "[Screen cleared]"
		},
	})
	r.Register(Command{
		Name:        "status",
		Description: "Show active agent and workspace status",
		Execute: func(args []string) string {
			return "Tercode Core: Running\nWorkspace: Active"
		},
	})
}

// Register adds a command.
func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name] = cmd
}

// Execute parses and runs a slash command string.
func (r *Registry) Execute(input string) (string, bool) {
	if !strings.HasPrefix(input, "/") {
		return "", false
	}
	parts := strings.Fields(input[1:])
	if len(parts) == 0 {
		return "", false
	}
	name := parts[0]
	args := parts[1:]

	cmd, ok := r.commands[name]
	if !ok {
		return fmt.Sprintf("Unknown command /%s. Type /help for available commands.", name), true
	}
	return cmd.Execute(args), true
}
