package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"github.com/librant/learn/k8s/operator/pod-controller/pkg/controller"
	"github.com/librant/learn/k8s/operator/pod-controller/pkg/signals"
)

// kubeConfig kube-config path
var kubeConfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "root-cmd",
	Short: "controller-controller damon",
	Long:  `controller-controller is a simple k8s controller which watch controller change`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Infof("controller-controller kubeconfig: %s", kubeConfig)
		stopCh := signals.SetupSignalHandler()
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			klog.Fatalln(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			klog.Fatalln(err)
		}

		podFactory := informers.NewSharedInformerFactory(clientset, time.Minute)
		podInformer := podFactory.Core().V1().Pods()
		podController := controller.New(clientset, podInformer)

		// 启动 informer
		go podFactory.Start(stopCh)

		if err := podController.Run(maxWorkNum, stopCh); err != nil {
			klog.Fatalln(err)
		}
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
		runtime.HandleError(err)
		os.Exit(1)
	}
}
