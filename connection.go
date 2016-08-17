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

  return Connection{disqueConn, cycle, nodes}, nil
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

func (c *Connection) Fetch(count int, timeout time.Duration, queues ...string) (Job, error){
  arguments := redis.Args{}.
    Add("TIMEOUT").Add(int64(timeout.Seconds() * 1000)).
    Add("COUNT").Add(count).
    Add("FROM").AddFlat(queues)

  values, err := redis.Values(c.Do("GETJOB", arguments...)); if err != nil {
    return Job{}, err
  }

  for _, value := range values {
    jobData, err := redis.Values(value, nil); if err != nil {
      return Job{}, err
    }

    if len(jobData) < 3 {
      return Job{}, errors.New("Malformed job fetched from Disque")
    }

    return Job{
      Queue:    string(jobData[0].([]byte)),
      ID:       string(jobData[1].([]byte)),
      Payload:  jobData[2].([]byte),
    }, nil
  }

  return Job{}, errors.New("timeout reached")
}
