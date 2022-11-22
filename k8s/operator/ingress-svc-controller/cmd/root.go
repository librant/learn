package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"github.com/librant/learn/k8s/operator/ingress-svc-controller/pkg/signals"
)

// kubeConfig kube-config path
var kubeConfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "root-cmd",
	Short: "ingress-svc-controller daemon",
	Long:  `ingress-svc-controller is a simple k8s ingress controller which watch service change`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Infof("controller-controller kubeconfig: %s", kubeConfig)
		// 1、获取 rest config 对象
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			klog.Fatalln(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			klog.Fatalln(err)
		}
		factory := informers.NewSharedInformerFactory(clientset, 0)
		serviceInformer := factory.Core().V1().Services()
		ingressInformer := factory.Networking().V1().Ingresses()

		stopCh := signals.SetupSignalHandler()
		go factory.Start(stopCh)

		factory.WaitForCacheSync(stopCh)

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
