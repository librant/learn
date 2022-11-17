package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/librant/learn/k8s/operator/pod-controller/pkg/controller"
)

// kubeConfig kube-config path
var kubeConfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "root-cmd",
	Short: "controller-controller damon",
	Long:  `controller-controller is a simple k8s controller which watch controller change`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("controller-controller kubeconfig: %s", kubeConfig)
		stopCh := make(chan struct{})
		podController := controller.New(kubeConfig)
		podController.Run(stopCh)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&kubeConfig, "config", "",
		"config file (default is ./.kube/kubeconfig)")

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if len(kubeConfig) != 0 {
		// Use config file from the flag.
		viper.SetConfigFile(kubeConfig)
		return
	}
	kubeConfig = defaultKubeConfigPath
}

// Execute cmd execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("rootCmd execute failed: %v", err)
		os.Exit(1)
	}
}
