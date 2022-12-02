```shell
pkg/cri/server/sandbox_run.go
```

在 Kubernetes 中定义的 Sandbox 接口：
```go
// PodSandboxManager contains methods for operating on PodSandboxes. The methods
// are thread-safe.
type PodSandboxManager interface {
	// RunPodSandbox creates and starts a pod-level sandbox. Runtimes should ensure
	// the sandbox is in ready state.
	RunPodSandbox(ctx context.Context, config *runtimeapi.PodSandboxConfig, runtimeHandler string) (string, error)
	// StopPodSandbox stops the sandbox. If there are any running containers in the
	// sandbox, they should be force terminated.
	StopPodSandbox(pctx context.Context, odSandboxID string) error
	// RemovePodSandbox removes the sandbox. If there are running containers in the
	// sandbox, they should be forcibly removed.
	RemovePodSandbox(ctx context.Context, podSandboxID string) error
	// PodSandboxStatus returns the Status of the PodSandbox.
	PodSandboxStatus(ctx context.Context, podSandboxID string, verbose bool) (*runtimeapi.PodSandboxStatusResponse, error)
	// ListPodSandbox returns a list of Sandbox.
	ListPodSandbox(ctx context.Context, filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error)
	// PortForward prepares a streaming endpoint to forward ports from a PodSandbox, and returns the address.
	PortForward(context.Context, *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error)
}
```

通过 grpc 的方式调用 containerd 中的 cri 接口：
```go
func (c *criService) RunPodSandbox(ctx context.Context, 
	r *runtime.RunPodSandboxRequest) (_ *runtime.RunPodSandboxResponse, retErr error) {}
```

这里针对非 Host 网络创建 pod 网络：
```shell
var netnsMountDir = "/var/run/netns"
```
// 创建 network namespace
--> sandbox.NetNS, err = netns.NewNetNS(netnsMountDir)
--> sandbox.NetNSPath = sandbox.NetNS.GetPath()
--> c.setupPodNetwork(ctx, &sandbox)

```go
// setupPodNetwork setups up the network for a pod
func (c *criService) setupPodNetwork(ctx context.Context, sandbox *sandboxstore.Sandbox) error {}
```
--> netPlugin = c.getNetworkPlugin(sandbox.RuntimeHandler)
```shell
kind: RuntimeClass
```

在 containerd 中的 criService 结构体中：
```go
// criService implements CRIService.
type criService struct {}
```
的成员变量, 根据接口中传入的 runtimeClass 的名称，获取指定的 CNI 插件实现：
```go
	// netPlugin is used to setup and teardown network when run/stop pod sandbox.
	netPlugin map[string]cni.CNI
```

可以继续追踪下，这里的 netPlugin 是在哪里进行初始化的；
```shell
pkg/cri/server/service_linux.go
```
```go
func (c *criService) initPlatform() (err error) {}
```
```go
	pluginDirs := map[string]string{
        defaultNetworkPlugin: c.config.NetworkPluginConfDir,
        }

	c.netPlugin = make(map[string]cni.CNI)
	for name, dir := range pluginDirs {
		max := c.config.NetworkPluginMaxConfNum
		if name != defaultNetworkPlugin {
			if m := c.config.Runtimes[name].NetworkPluginMaxConfNum; m != 0 {
				max = m
			}
		}
		// Pod needs to attach to at least loopback network and a non host network,
		// hence networkAttachCount is 2. If there are more network configs the
		// pod will be attached to all the networks but we will only use the ip
		// of the default network interface as the pod IP.
		i, err := cni.New(cni.WithMinNetworkCount(networkAttachCount),
			cni.WithPluginConfDir(dir),
			cni.WithPluginMaxConfNum(max),
			cni.WithPluginDir([]string{c.config.NetworkPluginBinDir}))
		if err != nil {
			return fmt.Errorf("failed to initialize cni: %w", err)
		}
		c.netPlugin[name] = i
	}
```
