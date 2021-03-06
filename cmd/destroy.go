package cmd

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/destroy"
	"github.com/craigmonson/colonize/prep"
)

type DestroyFlags struct {
	Environment string
	SkipRemote  bool
}

var destroyFlags = DestroyFlags{}

var yesToDestroy bool

const WARNING_MSG string = `All managed infrastructure will be deleted.
There is no undo. Only entering 'yes' will confirm this operation.

Do you wish to proceed with destroy: `

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy all defined resources in an environment",
	Long: `
This command will perform a "terraform destroy" command on your project for the 
specified environment. In effect, this will destroy all managed resources in the
given leaf or branch that the destroy command is run under.

# Example usage to destroy the "dev" environment: 
$ colonize destroy -e dev

# Example usage to destroy the "dev" environment, say yes to prompt
$ colonize destroy -e dev -y
`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(destroyFlags.Environment)
		if err != nil {
			CompleteFail(err.Error())
		}

		if !yesToDestroy {
			scan := bufio.NewScanner(os.Stdin)
			Log.Print(WARNING_MSG)
			scan.Scan()

			if scan.Text() != "yes" {
				CompleteFail("Destroy operation cancelled by user")
			}
		}

		err = Run("PREP", prep.Run, conf, Log, false, nil)
		if err != nil {
			CompleteFail("Prep failed to run: " + err.Error())
		}

		err = Run("DESTROY", destroy.Run, conf, Log, true, destroy.RunArgs{
			SkipRemote: destroyFlags.SkipRemote,
		})

		if err != nil {
			CompleteFail("Destroy failed to run: " + err.Error())
		}

		CompleteSucceed()
	},
}

func init() {
	addEnvironmentFlag(destroyCmd, &destroyFlags.Environment)
	addSkipRemoteFlag(destroyCmd, &destroyFlags.SkipRemote)
	destroyCmd.Flags().BoolVarP(
		&yesToDestroy,
		"accept",
		"y",
		false,
		"bypass 'accept' prompt by automatically accepting the destruction",
	)
	RootCmd.AddCommand(destroyCmd)
}
