
1、IPVS 的工作原理  
IPVS 是内核实现的四层负载均衡，是 LVS 负载均衡模块的实现；

三种负载均衡模式：
- Direct Routing
  - 工作在 L2 层，即通过 MAC 地址做的 LB，而非 IP
  - 不支持端口映射
- Tunneling (ipip)
  - 利用 IP 包封装 IP 包
  - 不支持端口映射
- NAT（Masq）
  - 支持端口映射，回程报文需要经过 IPVS Director

2、IPVS 的参数模式
--proxy-mode=ipvs 
--ipvs-shceduler:
- rr: 轮询
- lc: 最小连接数
- dh: 目的地址 hash
- sh：源地址 hash
- sed：最短延时

确保内核安装了 ipvs 模块（lsmod | grep ip_vs）

3、IPVS 模式实现原理

一旦创建 service 和 endpoints， kube-proxy 会做以下三件事：  
1）确保一块 dummy 网卡（kube-ipvs0）存在  
  创建 dummy 网卡，因为 IPVS 的 netfilter 钩子挂载 INPUT 链，
  我们需要把 Service 的访问 IP 绑定在 dummy 网卡上，让内核"觉得"虚 IP 就是本机 IP, 进而进入 INPUT 链
2）把 Service 的访问 IP 绑定在 dummy 网卡上
3）通过 socket 调用，创建 IPVS 的 virtual server 和 real server，
  分别对应 kubernetes 的 Service 和 Endpoints

4、conntrack：连接跟踪模块  
netfilter 框架实现的连续跟踪模块称作 conntrack  
在 DNAT 的过程中，conntrack 使用状态机启动并跟踪连接状态  
- NEW: 匹配连接的第一个数据包，发生在 SYN 数据包时
- ESTABLISHED: 匹配连接的响应及后续的包
- RELATED: 当一个连接与另一个 ESTABLISHED 状态的链接有关时
- INVALID: 匹配那些无法识别或没有任何状态的数据包

