
#### ingress-svc-controller 功能说明：

1、创建 service 资源，添加 annotation: "ingress/http: true" 来标识是否可以访问该服务；

2、当 controller 监听到 service 资源变动时：
- 新增：
    包含指定 annotation 时，创建 ingress 资源对象；
    不包含指定 annotation 时：忽略；
- 删除：
    删除 ingress 资源对象
- 更新：
    包含指定 annotation 时，检测 ingress 资源对象是否存在，不存在，则进行创建；
    不包含指定 annotation 时，检测 ingress 资源对象是否存在，存在，则删除；

