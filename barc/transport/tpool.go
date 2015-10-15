package transport
import (
	"github.com/vireshas/minimal_vitess_pool/pools"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/akaspin/bar/proto/wire"
	"time"
	"math/rand"
)

// Bar thrift client wrapper
type TClient struct {
	*wire.BarClient
	pool *TPool
}

func (c *TClient) Connect() (err error) {
	if !c.Transport.IsOpen() {
		err = c.Transport.Open()
	}
	return
}

func (c *TClient) Release() {
	if !c.Transport.IsOpen() {
		c.pool.Put(nil)
		return
	}
	c.pool.Put(c)
}

func (c *TClient) Close() {
	if c.BarClient.Transport.IsOpen() {
		c.BarClient.Transport.Close()
	}
}

type TPool struct {
	endpoints []string
	endpointsCap int
	ttl time.Duration

	transportFactory thrift.TTransportFactory
	protoFactory thrift.TProtocolFactory

	*pools.ResourcePool
}

func NewTPool(endpoints []string, bufferSize int, n int, ttl time.Duration) (res *TPool) {
	res = &TPool{
		endpoints: endpoints,
		endpointsCap: len(endpoints) - 1,
		ttl: ttl,
		transportFactory: thrift.NewTBufferedTransportFactory(bufferSize),
		protoFactory: thrift.NewTBinaryProtocolFactoryDefault(),
	}
	res.ResourcePool = pools.NewResourcePool(res.factory, n, n, ttl)
	return
}

func (p *TPool) Take() (res *TClient, err error) {
	r, err := p.Get(p.ttl)
	if err != nil {
		return
	}
	res = r.(*TClient)
	if err = res.Connect(); err != nil {
		p.Put(nil)
	}
	return
}

func (p *TPool) factory() (res pools.Resource, err error)  {
	// peek endpoint
	var endpoint string
	if p.endpointsCap == 0 {
		endpoint = p.endpoints[0]
	} else {
		endpoint = p.endpoints[rand.Intn(p.endpointsCap)]
	}

	var transport thrift.TTransport
	if transport, err = thrift.NewTSocket(endpoint); err != nil {
		return
	}
	transport = p.transportFactory.GetTransport(transport)

	client := wire.NewBarClientFactory(transport, p.protoFactory)
	res = &TClient{client, p}
	return
}
