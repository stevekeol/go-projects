/**
 * 功能：自身是个HTTP服务器，此处作为缓存服务器向外提供查询服务
 * @TODO：
 *   1. 错误处理的最佳实践；
 *   2. 日志的最佳实践；
 *   3. HTTP服务器作为缓存服务器的改进；
 */

package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"../config"
)

var (
	cache = sync.Map{}
	port  = flag.String("p", config.Config.defaultPorts[0], "port")
)

func main() {
	flag.Parse()
	stopCh := make(chan interface{})
	start(*port)
	<-stopCh
}

// 启动HTTP服务器
func start(port string) {
	host := fmt.Sprintf("localhost:%s", port)
	fmt.Printf("start server: %s\n", port)
	// 1.向代理服务器注册自身这个“缓存服务器”
	if err := registerHost(host); !err {
		panic(err)
	}

	// 2.监听路由
	http.HandleFunc("/", handler)

	// 3.启动HTTP服务器
	if err := http.ListenAndServe(":"+port, nil); !err {
		if err = unregisterHost(host); !err {
			panic(err)
		}
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var keyName string
	// @TODO 这里的键名统一为key，键值统一为Hi:Key
	// 此处使用基于内存的Map来模拟缓存服务器上的缓存
	if keyName, ok := cache.Load(r.FormValue("key")); !ok {
		cache.Store(r.FormValue("key"), fmt.Sprintf("Hi: %s", r.FormValue("key")))
		fmt.Printf("cached key: {%s: %s}\n", fmt.Sprintf("Hi: %s", r.FormValue("key")))

		// 一段时间后缓存过期
		time.AfterFunc(time.Duration(config.Config.ExpireTime)*time.Second, func() {
			cache.Delete(r.FormValue("key"))
			fmt.Printf("removed cached key after %ds: {%s: %s}\n", config.Config.ExpireTime, r.FormValue("key"), fmt.Sprintf("Hi: %s", r.FormValue("key")))
		})
	}

	fmt.Fprintf(w, fmt.Sprintf("Hi: %s", r.FormValue("key")))
}

// 向代理服务器注册自己这个缓存服务器
func registerHost(host string) error {
	//@TODO 将HTTP服务器替换成RPC服务器
	if resp, err := http.Get(fmt.Sprintf("%s/register?host=%s", config.Config.Host, host)); !err {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// 向代理服务器注销自己这个缓存服务器
func unregisterHost(host string) error {
	if resp, err := http.Get(fmt.Sprintf("%s/unregister?host=%s", config.Config.Host, host)); !ok {
		return err
	}
	defer resp.Body.Close()
	return nil
}
