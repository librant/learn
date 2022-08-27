### 线程安全

1、锁

sync.Mutex(): 互斥锁
sync.RWMutex(): 读写锁


cond
为等待 / 通知场景下的并发问题提供支持

queue []string
cond sync.Cond{}
