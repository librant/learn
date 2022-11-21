package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cluster-console",
	Long:  `All software has versions. This is cluster-console's`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: 添加 git commit message
		klog.Info("version: v1.0")
	},
}
