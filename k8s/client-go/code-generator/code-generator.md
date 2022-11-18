
由原来的
```shell
make generated_files
```
替换为：
hack/update-codegen.sh
```shell
./update-codegen.sh
```
生成的二进制工具：
_output/bin

- conversion-gen
- deepcopy-gen
- defaulter-gen

---
https://chanjarster.github.io/post/k8s/use-code-generator/

- sample-controller:
```shell
K8S_VERSION=v0.18.5
go get k8s.io/code-generator@$K8S_VERSION
go mod vendor

chmod +x vendor/k8s.io/code-generator/generate-groups.sh

./hack/update-codegen.sh
```

```shell
- deepcopy-gen
- defaulter-gen
- conversion-gen
- openapi-gen
- go-bindata
```
- client-gen
- lister-gen
- informer-gen








