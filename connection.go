package disco

import(
  "errors"
  "os"
  "strings"
  "time"

  "github.com/garyburd/redigo/redis"
)


type Connection struct {
  redis.Conn

  Cycle int
  Nodes []string

  // This is ugly and duplicated, but we need an explicit reference to
  // the connection so that we can leverage Redigo's Pool features.
  disqueConnection redis.Conn
}

func NewConnection(cycle int) (Connection, error){
  return NewConnectionToURLS(cycle, os.Getenv("DISQUE_NODES"))
}

func NewConnectionToURLS(cycle int, nodes string)  (Connection, error) {
  return NewConnectionToNodes(cycle, strings.Split(nodes, ",")...)
}

func NewConnectionToNodes(cycle int, nodes ...string) (Connection, error){
  disqueConn, err := connectToFirstAvailableNode(nodes...); if err != nil {
    return Connection{Cycle: cycle, Nodes: nodes}, err
  }

  return Connection{disqueConn, cycle, nodes, disqueConn}, nil
}

func connectToFirstAvailableNode(nodes ...string) (redis.Conn, error) {
  for _, node := range nodes {
    conn, err := redis.Dial("tcp", node); if err == nil {
      return conn, err
    }
  }

  return nil, errors.New("No available nodes")
}

func (c *Connection) AddJob(queue string, payload string, pushTimeout string) (string ,error) {
  timeout, err := time.ParseDuration(pushTimeout); if err != nil {
    return "", err
  }

  arguments := redis.Args{}.
    Add(queue).
    Add(payload).
    Add(int64(timeout.Seconds() * 1000))

  return redis.String(c.Do("ADDJOB", arguments...))
}
