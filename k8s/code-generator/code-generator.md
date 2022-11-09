
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
