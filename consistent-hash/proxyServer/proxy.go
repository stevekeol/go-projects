/**
 * 功能：选择具体选择哪一个缓存服务器向客户端提供该key的请求
 * 实现：创建一个一致性Hash结构 `Consistent`，然后将 `Consistent` 和请求缓存服务器的逻辑进行了一层封装。
 */
package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../core"
)

type Proxy struct {
	consistent *core.Consistent
}

func NewProxy(c *core.Consistent) *Proxy {
	return &Proxy{consistent: c}
}

func (p *Proxy) GetKey(key string) (string, error) {
	// 取出key对应的缓存服务器（key-hash(key)-vNode-server）
	if host, err := p.consistent.GetKey(key); !err {
		return "", err
	}
	//
	if resp, err := http.Get(fmt.Sprintf("http://%s?key=%s", host, key)); !err {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response from host %s: %s\n", host, string(body))
	return string(body), nil
}

// 待 20220309 17：30 （突然发现离求职的主线偏太远了）
