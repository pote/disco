package disco

import(
  "os"
  "strings"
  "time"

  "github.com/garyburd/redigo/redis"
)

type Pool struct {
  disqueConnections redis.Pool
  Cycle int
  Nodes []string
}

func NewPool(maxIdle, maxActive int, idleTimeout string, cycle int) (Pool, error) {
  return NewPoolToURLS(maxIdle, maxActive, idleTimeout, cycle, os.Getenv("DISQUE_NODES"))
}

func NewPoolToURLS(maxIdle, maxActive int, idleTimeout string, cycle int, urls string) (Pool, error) {
  return NewPoolToNodes(maxIdle, maxActive, idleTimeout, cycle, strings.Split(urls, ",")...)
}

func NewPoolToNodes(maxIdle, maxActive int, idleTimeout string, cycle int, nodes ...string) (Pool, error) {
  timeout, err := time.ParseDuration(idleTimeout); if err != nil {
    return Pool{}, err
  }

  disquePool := redis.Pool{
    MaxIdle: maxIdle,
    MaxActive: maxActive,
    IdleTimeout: timeout,
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
  _, err = c.Do("PING")

  p := Pool{
    disqueConnections: disquePool,
    Cycle: cycle,
    Nodes: nodes,
  }

  return p, err
}

func (p *Pool) Get() Connection {
  c := p.disqueConnections.Get()

  return Connection{c, p.Cycle, p.Nodes}
}

func (p *Pool) NewFunnel(queues ...string) Funnel {
  return NewFunnel(p, queues...)
}
