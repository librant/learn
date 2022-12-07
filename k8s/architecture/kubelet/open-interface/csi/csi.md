
#### CSI（Container Storage Interface）

1) CSI 持久化卷具有以下字段可供用户指定   
- driver：一个字符串值，指定要使用卷驱动程序的名称
- volumeHandle：一个字符串值，唯一标识从 CSI 卷插件的 CreateVolume 调用返回的卷名
- readOnly：一个可选的布尔值，指示卷是否被发布为只读。默认是 false

2) 动态配置
- StorageClass：支持动态配置的 CSI Storage 插件启用自动创建

3) 附着和挂载   
- 在任何的 pod 或者 pod 的 template 中引用绑定到 CSI volume 上的 PVC 

4) 创建 CSI 驱动    
提供以下 sidecar 容器部署方案：   
- External-attacher
- External-provisioner
- Cluster Driver Registrar
- Node Driver Registrar



