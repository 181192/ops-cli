package wrapper

import (
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/util"
	"github.com/spf13/cobra"
)

// Command will create the `wrapper` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	wrappers := MakeWrappers()
	helpText := "Wrapper commands"
	helpTextLong := "Available wrapper commands:\n\n"

	for _, w := range wrappers {
		helpTextLong += fmt.Sprintf("  - %s\t%s\n", w.Name, w.Description)
	}

	wrapperCmd := &cobra.Command{
		Use:     "wrapper",
		Aliases: []string{"w"},
		Short:   helpText,
		Long:    helpTextLong,
	}

	return generateWrapperCommands(flagGrouping, wrapperCmd, wrappers)
}

func generateWrapperCommands(flagGrouping *cmdutils.FlagGrouping, wrapperCmd *cobra.Command, wrappers []Wrapper) *cobra.Command {
	for _, wrapper := range wrappers {
		cmdutils.AddResourceCmd(flagGrouping, wrapperCmd, func(cmd *cmdutils.Cmd) {
			cmd.CobraCommand.Use = wrapper.Name
			cmd.CobraCommand.Run = func(_ *cobra.Command, args []string) {
				for _, w := range wrappers {
					if w.Name == cmd.CobraCommand.Use {
						wrapper = w
						break
					}
				}

				util.ExecuteCmd(cmd.CobraCommand, wrapper.Executable, args)
			}
		})
	}
	return wrapperCmd
}
