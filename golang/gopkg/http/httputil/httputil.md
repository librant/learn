#### httputil 

#### 1、httputil.ReverseProxy 功能

- 自定义修改响应包体
- 连接池
- 错误信息自定义处理
- 支持 websocket 服务
- 自定义负载均衡
- https 代理
- URL 重写

#### 2、httputil.ReverseProxy 结构
- Director
    - 当接收到客户端请求时，ServeHTTP 函数首先调用 Director 函数对接受到的请求体进行修改，例如修改请求的目标地址、请求头等；然后使用修改后的请求体发起新的请求
- ModifyResponse
    - 接收到响应后，调用 ModifyResponse 函数对响应进行修改，最后将修改后的响应体拷贝并响应给客户端，这样就实现了反向代理的整个流程
    - 用于修改响应结果及 HTTP 状态码，当返回结果 error 不为空时，会调用 ReverseProxy.ErrorHandler
- ErrorHandler
    - 用于处理后端和 ModifyResponse 成员返回的错误信息，默认将返回传递过来的错误信息，并返回 HTTP 502错误码    

---
参考文档：
1、[Golang ReverseProxy 分析](https://pandaychen.github.io/2021/07/01/GOLANG-REVERSEPROXY-LIB-ANALYSIS/)