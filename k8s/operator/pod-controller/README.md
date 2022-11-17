### 主要功能：实现 pod-controller

1、创建CRD（Custom Resource Definition），令k8s明白我们自定义的API对象；

2、编写代码，将CRD的情况写入对应的代码中，然后通过自动代码生成工具，将controller之外的informer，client等内容较为固定的代码通过工具生成；

3、编写controller，在里面判断实际情况是否达到了API对象的声明情况，如果未达到，就要进行实际业务处理，而这也是controller的通用做法；