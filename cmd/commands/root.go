package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-tech/cmd/commands/api"
	"go-tech/internal/logging"
	"io"
	"os"
)

func newRootCmd(out io.Writer, args []string) (*cobra.Command, error) {

	var cmd = &cobra.Command{
		Use:   "go-tech",
		Short: "the model manage for iqi business",
	}
	flags := cmd.PersistentFlags()
	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.Parse(args)
	api.AddApi(cmd)
	return cmd, nil

}
func Execute() {

	fmt.Println(os.Args)
	cmd, err := newRootCmd(os.Stdout, os.Args[1:])
	if err != nil {
		logging.ErrorStack(nil, err).Send()
		os.Exit(-1)
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
