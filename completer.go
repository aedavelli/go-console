package console

import (
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func getSuggestions(d prompt.Document, m cmdMap) []prompt.Suggest {

	s := []prompt.Suggest{}

	tok, err := parseCommandLine(d.Text)
	if err != nil || len(tok) == 0 {
		return s
	}

	if len(tok) == 1 {
		//Iterate over top level commands
		for _, lm := range []cmdMap{context[""], m} {
			for n, cmd := range lm {
				if n == cmd.Name() {
					s = cmdSuggestions(cmd, s, true)
				}
			}
		}
		return s
	}

	for _, lm := range []cmdMap{context[""], m} {
		if cmd, ok := lm[tok[0]]; ok {
			if c := findLastValidCommand(cmd, tok[1:]); c != nil {
				if c.HasSubCommands() {
					s = cmdSuggestions(c, s, false)
				} else {
					s = cmdSuggestions(c, s, true)
				}
			}
			return s
		}
	}
	return s
}

func cmdSuggestions(c *cobra.Command, s []prompt.Suggest, this bool) []prompt.Suggest {
	var cmds []*cobra.Command
	if this {
		cmds = make([]*cobra.Command, 0)
		cmds = append(cmds, c)
	} else {
		cmds = c.Commands()
	}

	for _, c := range cmds {
		// Here we've matching subcommand
		s = append(s, prompt.Suggest{Text: c.Name(), Description: c.Short})
		// Look for aliases
		for _, a := range c.Aliases {
			s = append(s, prompt.Suggest{Text: a, Description: c.Short})
		}

		if this {
			// Populate flags
			if fs := c.Flags(); fs != nil {
				fs.VisitAll(func(f *pflag.Flag) {
					if f.Shorthand != "" {
						s = append(s, prompt.Suggest{Text: "-" + f.Shorthand, Description: f.Usage})
					}
					if f.Name != "" {
						s = append(s, prompt.Suggest{Text: "--" + f.Name, Description: f.Usage})
					}
				})
			}

			// Check whether any argument completer available for this command
			// If this is global command context should be zero string value
			ctx := ""
			if _, ok := context[""][c.Name()]; !ok {
				ctx = presentCtx
			}

			if f, ok := argCompleter[cmdContextPath(c, ctx)]; ok {
				for key, val := range f() {
					s = append(s, prompt.Suggest{Text: key, Description: val})
				}
			}
		}
	}
	return s
}

func findLastValidCommand(c *cobra.Command, args []string) *cobra.Command {
	if !c.HasSubCommands() || len(args) < 2 {
		return c
	}

	if cmd := findNext(c, args[0]); cmd != nil {
		return findLastValidCommand(cmd, args[1:])
	}
	return nil
}

func findNext(c *cobra.Command, next string) *cobra.Command {
	for _, cmd := range c.Commands() {
		if cmd.Name() == next || cmd.HasAlias(next) {
			return cmd
		}
	}
	return nil
}

func Completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	w := d.GetWordBeforeCursor()
	if w == "" {
		return s
	}

	if m, ok := context[presentCtx]; ok && presentCtx != "" {
		s = getSuggestions(d, m)
	} else {
		s = getSuggestions(d, nil)
	}

	return prompt.FilterHasPrefix(s, w, true)
}
