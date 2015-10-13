package transport
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"net/rpc"
	"net/url"
	"time"
)

type RPCClient struct  {
	endpoint *url.URL
	pool *RPCPool
	client *rpc.Client
}

func (c *RPCClient) Connect() (err error) {
	if c.client == nil {
		c.client, err = rpc.DialHTTPPath(
			"tcp",
			c.endpoint.Host,
			c.endpoint.Path + "/rpc")
	}
	return
}

func (c *RPCClient) Call(serviceMethod string, args interface{}, reply interface{}) (err error) {
	if err = c.Connect(); err != nil {
		return
	}
	err = c.client.Call(serviceMethod, args, reply)
	return
}

func (c *RPCClient) Release() {
	if c.client == nil {
		c.pool.Put(nil)
		return
	}
	c.pool.Put(c)
}

func (c *RPCClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
	c.client = nil
}


type RPCPool struct {
	endpoint string
	rpcEndpoints string
	timeout time.Duration
	*pools.ResourcePool
}

func NewRPCPool(n int, ttl time.Duration,
	endpoint string, rpcEndpoints []string) (res *RPCPool) {
	res = &RPCPool{
		endpoint: endpoint,
		timeout: ttl,
	}
	res.ResourcePool = pools.NewResourcePool(res.factory, n, n, ttl)
	return
}

func (p *RPCPool) Take() (res *RPCClient, err error) {
	r1, err := p.Get(p.timeout)
	if err != nil {
		return
	}
	res = r1.(*RPCClient)
	// connect
	if err = res.Connect(); err != nil {
		res.Close()
		res.Release()
		res = nil
	}
	return
}

func (p *RPCPool) factory() (res pools.Resource, err error) {
	u, err := url.Parse(p.endpoint)
	if err != nil {
		return
	}
	res = &RPCClient{endpoint: u, pool: p}
	return
}


