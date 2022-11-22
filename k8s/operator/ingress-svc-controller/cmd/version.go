package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is ingress-svc-controller's`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: 添加 git commit message
		klog.Info("version: v1.0")
	},
}
