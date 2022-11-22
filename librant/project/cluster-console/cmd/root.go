package cmd

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"

	"github.com/librant/learn/librant/project/cluster-console/internal/handler"
	"github.com/librant/learn/librant/project/cluster-console/internal/pkg/config"
)

// kubeConfig kube-config path
var kubeConfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "root-cmd",
	Short: "cluster-console for container login",
	Long:  `cluster-console is a simple k8s cluster which for container login`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Infof("cluster-console run begin...")
		if err := config.Init(kubeConfig); err != nil {
			klog.Fatal(err)
		}
		router := gin.Default()
		// 静态资源加载，本例为 css,js 以及资源图片
		router.StaticFS("/static/", http.Dir("."))
		router.StaticFile("/favicon.ico", "/static/favicon.ico")

		router.GET("/login", handler.IndexHandler)
		klog.Fatal(router.Run(":8080"))
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
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		klog.Errorf("cluster-console read cluster config failed: %v", err)
	}
	kubeConfig = viper.ConfigFileUsed()
}

// Execute cmd execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		runtime.HandleError(err)
		os.Exit(1)
	}
}
