package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/librant/learn/k8s/operator/service-ingress-controller/cmd/pkg"
)

var (
	runCommand = "service-ingress-controller"

	// Used for flags.
	kubeconfig = "./.kube/kubeconfig"
)

var (
	rootCmd = &cobra.Command{
		Use:   runCommand,
		Short: "run is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Printf("kubeconfig: %s", kubeconfig)

			// 1 指定 k8s 的配置文件
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				log.Panicln(err)
			}
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				log.Panicln(err)
			}
			factory := informers.NewSharedInformerFactory(clientset, 0)

			serviceInformer := factory.Core().V1().Services()
			ingressInformer := factory.Networking().V1().Ingresses()

			stopCh := make(chan struct{})
			factory.Start(stopCh)
			factory.WaitForCacheSync(stopCh)

			c := pkg.New(clientset, serviceInformer, ingressInformer)
			c.Run(stopCh)
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
