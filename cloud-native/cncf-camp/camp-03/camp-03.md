#### Docker 核心技术

##### 关于 namespace 的常用操作

- 查看当前系统的 namespace:
lsns -t <type>

- 查看某进程的 namespace:
ls -la /proc/<pid>/ns/

- 进入某 namespace 运行命令：
nssenter -t <pid> -n ip addr

##### 关于 cgroup 的常用操作

