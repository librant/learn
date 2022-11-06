package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	runCommand = "rest-client-cobra"
)

var (
	// Used for flags.
	kubeconfig = "./.kube/kubeconfig"

	rootCmd = &cobra.Command{
		Use:   runCommand,
		Short: "run is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Printf("run here")
			log.Printf("kubeconfig: %s", kubeconfig)

			// 1 指定 k8s 的配置文件
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				log.Panicln(err)
			}
			// 2 配置 API 路径
			config.APIPath = "api"

			// 3 设置分组版本
			config.GroupVersion = &corev1.SchemeGroupVersion

			// 4 配置数据的编解码器
			config.NegotiatedSerializer = scheme.Codecs

			// 5 实例化 rest client
			restClient, err := rest.RESTClientFor(config)
			if err != nil {
				log.Panicln(err)
			}

			// 6 定义返回接收值
			result := &corev1.PodList{}
			if err := restClient.Get().
				Namespace("kube-system").
				Resource("pods").
				VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
				Do(context.Background()).
				Into(result); err != nil {
				log.Panicln(err)
			}

			for _, item := range result.Items {
				log.Printf("namespace: %s name: %s\n", item.Namespace, item.Name)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&kubeconfig, "config", "./.kube/kubeconfig", "config file (default is ./.kube/kubeconfig")
	rootCmd.PersistentFlags().StringP("author", "a", "librant", "author name for copyright attribution")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "librant <librant@aliyun.com>")

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if kubeconfig != "" {
		// Use config file from the flag.
		viper.SetConfigFile(kubeconfig)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
