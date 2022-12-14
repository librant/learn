package app

import (
	"flag"
	"log"

	"github.com/spf13/cobra"

	"github.com/librant/learn/k8s/client-go/demon/cobra-rest-client/app/options"
)

var (
	// the command name
	cmdName = "run"
)


// NewRestClientCommand 生成 rest client 命令
func NewRestClientCommand() *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  cmdName,
		Short: "rest client run",
		Long: "ignore",
		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("helloworld: %v", args)
			return nil
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine) // for --boot-id-file and --machine-id-file

	_ = cmd.MarkFlagFilename("config", "yaml", "yml", "json")

	return cmd
}
