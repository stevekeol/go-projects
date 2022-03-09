# ProxyServer

代理服务器：来选择具体选择哪一个缓存服务器向客户端提供该key的请求

其逻辑很简单：就是创建一个一致性Hash结构 `Consistent`，然后将 `Consistent` 和请求缓存服务器的逻辑进行了一层封装。