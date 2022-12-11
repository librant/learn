
#### Service Mesh 

1) Service Mesh 特点
- 应用程序间通信的中间层
- 轻量级网络代理
- 应用程序无感知
- 解耦应用程序的重试，超时，监控，追踪，服务发现

2) Service Mesh 基本原理
负责服务之间的调用，限流，熔断和监控    

- Service Mesh 的架构
![img.png](img.png)

3) Service Mesh 方案    
- Linkerd 
- Istio

1、Istio    
- 架构图
![img_2.png](img_2.png)

---
1) 服务网格
- 流量治理
  - 负载均衡策略
  - 按比例切流
  - 按自定义业务策略切流
  - 错误重试
  - 熔断保护
- 可观测性
  - 请求参数监控 (延时、5XX 分析统计)
  - 开箱即用调用链
- 安全性
  - mtls 通道级安全
  - jwt token 安全认证







