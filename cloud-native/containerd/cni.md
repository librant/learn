```shell
github.com/containerd/go-cni/cni.go
```

```go
type CNI interface {
	// Setup setup the network for the namespace
	Setup(ctx context.Context, id string, path string, opts ...NamespaceOpts) (*Result, error)
	// SetupSerially sets up each of the network interfaces for the namespace in serial
	SetupSerially(ctx context.Context, id string, path string, opts ...NamespaceOpts) (*Result, error)
	// Remove tears down the network of the namespace.
	Remove(ctx context.Context, id string, path string, opts ...NamespaceOpts) error
	// Check checks if the network is still in desired state
	Check(ctx context.Context, id string, path string, opts ...NamespaceOpts) error
	// Load loads the cni network config
	Load(opts ...Opt) error
	// Status checks the status of the cni initialization
	Status() error
	// GetConfig returns a copy of the CNI plugin configurations as parsed by CNI
	GetConfig() *ConfigResult
}
```