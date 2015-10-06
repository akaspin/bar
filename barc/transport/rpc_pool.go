package transport
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"net/rpc"
	"net/url"
	"time"
)

type RPCClient struct  {
	endpoint *url.URL
	pool *rpcClientPool
	c *rpc.Client
}

func (c *RPCClient) Connect() (err error) {
	if c.c == nil {
		c.c, err = rpc.DialHTTPPath("tcp", c.endpoint.Host, c.endpoint.Path + "/rpc")
	}
	return
}

func (c *RPCClient) Call(serviceMethod string, args interface{}, reply interface{}) (err error) {
	if err = c.Connect(); err != nil {
		return
	}
	err = c.c.Call(serviceMethod, args, reply)
	return
}

func (c *RPCClient) Release() {
	if c.c == nil {
		c.pool.p.Put(nil)
		return
	}
	c.pool.p.Put(c)
}

func (c *RPCClient) Close() {
	if c.c != nil {
		c.c.Close()
	}
	c.c = nil
}

type rpcClientPool struct {
	endpoint *url.URL
	p *pools.ResourcePool
}

func newRpcClientPool(endpoint *url.URL, n int, timeout time.Duration) (res *rpcClientPool) {
	res = &rpcClientPool{endpoint: endpoint}
	res.p = pools.NewResourcePool(res.factory, n, n, timeout)
	return
}

func (p *rpcClientPool) take() (res *RPCClient, err error) {
	c, err := p.p.TryGet()
	if err != nil {
		return
	}
	res = c.(*RPCClient)
	err = res.Connect()
	return
}

func (p *rpcClientPool) factory() (res pools.Resource, err error) {
	res = &RPCClient{endpoint: p.endpoint, pool: p}
	return
}


type RPCPool struct {
	n int
	timeout time.Duration
	pools map[string]*rpcClientPool
}

func NewRPCPool(n int, timeout time.Duration) *RPCPool {
	return &RPCPool{
		n, timeout,
		map[string]*rpcClientPool{},
	}
}

func (p *RPCPool) Take(endpoint string) (res *RPCClient, err error) {
	pool, exists := p.pools[endpoint]
	if exists {
		res, err = pool.take()
		return
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return
	}
	p.pools[endpoint] = newRpcClientPool(u, p.n, p.timeout)
	res, err = p.pools[endpoint].take()
	return
}

func (p *RPCPool) Close() {
	for _, pool := range p.pools {
		pool.p.Close()
	}
}

