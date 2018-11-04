package console

import (
	"strings"

	"github.com/spf13/cobra"
)

func RegisterCommand(c *cobra.Command, f ArgCompleter) {
	RegisterCommandWithCtx(c, "", f)
}

func RegisterCommandWithCtx(c *cobra.Command, ctx string, f ArgCompleter) {
	cm, ok := context[ctx]
	if !ok {
		cm = make(cmdMap)
		context[ctx] = cm
	}

	cm[c.Name()] = c
	// Add aliases for command map
	for _, alias := range c.Aliases {
		cm[alias] = c
	}

	// Register the argument completer
	if f != nil {
		argCompleter[cmdContextPath(c, ctx)] = f
	}

}

func cmdContextPath(c *cobra.Command, ctx string) string {
	fc := ctx + "/" + c.CommandPath()
	return strings.Replace(fc, " ", "/", -1)
}
