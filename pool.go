package disco

import(
  "os"
  "strings"
  "time"

  "github.com/garyburd/redigo/redis"
)

type Pool struct {
  Connections redis.Pool
  Cycle int
  Nodes []string
}

func NewPool(maxIdle, maxActive, cycle int, idleTimeout time.Duration) (Pool, error) {
  return NewPoolToURLS(maxIdle, maxActive, cycle, idleTimeout, os.Getenv("DISQUE_NODES"))
}

func NewPoolToURLS(maxIdle, maxActive, cycle int, idleTimeout time.Duration, urls string) (Pool, error) {
  return NewPoolToNodes(maxIdle, maxActive, cycle, idleTimeout, strings.Split(urls, ",")...)
}

func NewPoolToNodes(maxIdle, maxActive, cycle int, idleTimeout time.Duration, nodes ...string) (Pool, error) {
  disquePool := redis.Pool{
    MaxIdle: maxIdle,
    MaxActive: maxActive,
    IdleTimeout: idleTimeout,
    Dial: func () (redis.Conn, error) {
      return connectToFirstAvailableNode(nodes...)
    },
    TestOnBorrow: func(c redis.Conn, t time.Time) error {
      _, err := c.Do("PING")
      return err
    },
  }

  c := disquePool.Get()
  defer c.Close()

  _, err := c.Do("PING")

  p := Pool{
    Connections: disquePool,
    Cycle: cycle,
    Nodes: nodes,
  }

  return p, err
}

func (p *Pool) Get() Connection {
  c := p.Connections.Get()

  return Connection{c, p.Cycle, p.Nodes}
}

func (p *Pool) NewFunnel(queues ...string) Funnel {
  return NewFunnel(p, 1, time.Second * 100, queues...)
}

func (p *Pool) NewFunnelWithOptions(fetchCount int, fetchTimeout time.Duration, queues ...string) Funnel {
  return NewFunnel(p, fetchCount, fetchTimeout, queues...)
}
