package console

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func listCmdsFromCtx(m cmdMap) {
	for n, cmd := range m {
		// Skip aliases
		if n == cmd.Name() {
			fmt.Println("\t", cmd.NameAndAliases())
		}
	}
}

func listCmds() {
	fmt.Println("Available commands")

	// List the commands at all contexts
	if m, ok := context[""]; ok {
		listCmdsFromCtx(m)
	}

	// List the commands at this context
	if m, ok := context[presentCtx]; ok && presentCtx != "" {
		listCmdsFromCtx(m)
	}
}

func init() {

	switchCmd := &cobra.Command{
		Use:   "switch [mode]",
		Short: "Change the console mode",
		Long:  "Change the console mode to one of the following \"admin\", \"config\"",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			presentCtx = args[0]
			return
		},
	}

	exitCmd := &cobra.Command{
		Use:     "exit",
		Short:   "Exit the console",
		Long:    "Exit the console",
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"quit"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Bye!")
			os.Exit(0)
			return
		},
	}

	RegisterCommand(exitCmd, nil)
	RegisterCommand(switchCmd, func() []string {
		s := make([]string, 0)
		for key, _ := range context {
			if key != "" {
				s = append(s, key)
			}
		}
		return s
	})

}

func handleCmd(args []string, cmds cmdMap) {
	// No need to verify args length. Here means it's already validated
	if cmd, ok := cmds[args[0]]; ok {
		resetAllCmdFlags(cmd)
		cmd.SetArgs(args[1:])
		cmd.Execute()
	} else {
		listCmds()
	}
}

func resetCmdFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	if fs != nil {
		fs.VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}
}

func resetAllCmdFlags(p *cobra.Command) {
	for _, cmd := range p.Commands() {
		resetCmdFlags(cmd)
	}
}

func Executor(in string) {
	args, _ := parseCommandLine(in)
	if len(args) == 0 {
		// No input provided. Provide list of commands available in this context and return
		listCmds()
		return
	}

	// Check whether the command is any of global commands exit, quit or switch
	cm := context[""]
	if cmd, ok := cm[args[0]]; ok {
		cmd.SetArgs(args[1:])
		cmd.Execute()
		return
	}

	if presentCtx != "" {
		if cm, ok := context[presentCtx]; ok {
			handleCmd(args, cm)
		} else {
			listCmds()
		}

	}
}
