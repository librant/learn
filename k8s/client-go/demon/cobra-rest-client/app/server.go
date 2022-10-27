package app

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/librant/learn/k8s/client-go/demon/cobra-rest-client/app/options"
)

// NewRestClientCommand 生成 rest client 命令
func NewRestClientCommand() *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  "rest-client",
		Long: "k8s rest client",
		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine) // for --boot-id-file and --machine-id-file

	_ = cmd.MarkFlagFilename("config", "yaml", "yml", "json")

	return cmd
}
