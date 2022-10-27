package options

import (
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
)

// Options 参数信息
type Options struct {
	// ConfigFile is the location of the scheduler server's configuration file.
	ConfigFile string

	// clientConnection specifies the kubeconfig file and client connection settings for the proxy
	// server to use when communicating with the apiserver.
	ClientConnection componentbaseconfig.ClientConnectionConfiguration

	Logs *logs.Options
	// Flags hold the parsed CLI flags.
	Flags *cliflag.NamedFlagSets
}

// NewOptions returns default scheduler app options.
func NewOptions() *Options {
	o := &Options{
		Logs: logs.NewOptions(),
	}

	o.initFlags()
	return o
}

// initFlags initializes flags by section name.
func (o *Options) initFlags() {
	if o.Flags != nil {
		return
	}

	nfs := cliflag.NamedFlagSets{}
	fs := nfs.FlagSet("misc")
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file.")
	logsapi.AddFlags(o.Logs, nfs.FlagSet("logs"))

	o.Flags = &nfs
}

// AddFlags adds flags to fs and binds them to options.
func (o *Options) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&o.ClientConnection.Kubeconfig, "kubeconfig", o.ClientConnection.Kubeconfig,
		"Path to kubeconfig file with authorization information (the master location can be overridden by the master flag).")

	return
}
