package disco

import(
  "time"
  "errors"
  "log"

  "github.com/garyburd/redigo/redis"
)

type Consumer struct {
  Queues []string
  IncomingJobs chan Job
  Connections *Pool
}

func (c *Consumer) Fetch(count int, fetchTimeout string) (Job, error){
  timeout, err := time.ParseDuration(fetchTimeout); if err != nil {
    return Job{}, err
  }

  arguments := redis.Args{}.
    Add("TIMEOUT").Add(int64(timeout.Seconds() * 1000)).
    Add("COUNT").Add(count).
    Add("FROM").AddFlat(c.Queues)

  connection := c.Connections.Get()
  defer connection.Close()

  values, err := redis.Values(connection.Do("GETJOB", arguments...)); if err != nil {
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


// This is a blocking call, you'll regularly want to execute it within a goroutine.
func (c *Consumer) FetchIntoChannel(count int, fetchTimeout string) {
  for {
    connection := c.Connections.Get()

    for {
      job, err := c.Fetch(count, fetchTimeout); if err != nil {
        log.Printf("Error fetching jobs in background: %v\n", err.Error())
        break
      }

      c.IncomingJobs <- job
    }

    connection.Close()
  }
}
