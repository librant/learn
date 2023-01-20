
###

```go
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
    JitterUntil(f, period, 0.0, true, stopCh)
}
```

```go
func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {}
```
- sliding: 确定系统周期调用的时候，计时的起点，如果是 TRUE 的话，就是 f 函数执行完成后计时周期。

Until(): 传入 true，Until是  f 函数运行开始前开始计时；
NonSlidingUntil(): 传入 false，f 函数运行结束后开始计时；
