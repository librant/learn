package app

import (
	"github.com/spf13/cobra"

	"github.com/librant/learn/k8s/client-go/demon/cobra-rest-client/app/options"
)

// NewRestClientCommand 生成 rest client 命令
func NewRestClientCommand() *cobra.Command {
	s := options.NewOptions()

	cmd := &cobra.Command{
		Use:  "rest-client",
		Long: "k8s rest client",
		// stop printing usage when the command errors
		SilenceUsage: true,
	}
	return cmd
}
