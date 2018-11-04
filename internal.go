package console

import (
	"github.com/spf13/cobra"
)

type ArgCompleter func() []string
type cmdMap map[string]*cobra.Command
type ctxMap map[string]cmdMap
type argMap map[string]ArgCompleter

var (
	presentCtx   = ""
	appName      = "console"
	context      = make(ctxMap, 0)
	argCompleter = make(argMap, 0)
)
