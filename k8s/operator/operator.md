
// 自动生成
operator-sdk init --domain=example.com --repo=github.com/librant/learn/k8s/operator/memcached-operator

operator-sdk create api --group cache --version v1 --kind Memcached --resource=true --controller=true

make docker-build IMG=librant/memcache:v1.0.0

// 参考资料
https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/