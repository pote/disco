package disco

import(
  "log"
  "time"
)

type Funnel struct {
  Queues []string
  Incoming chan Job
  Outgoing chan Job
  Connections *Pool

  FetchTimeout time.Duration
  FetchCount int

  Closed bool
}

func NewFunnel(pool *Pool, fetchCount int, fetchTimeout time.Duration, queues ...string) Funnel {
  incoming := make(chan Job)
  outgoing := make(chan Job)

  funnel := Funnel{
    Queues: queues,
    Incoming: incoming,
    Outgoing: outgoing,
    Connections: pool,
    FetchCount: fetchCount,
    FetchTimeout: fetchTimeout,
  }

  go funnel.Dispatch()
  go funnel.Listen()

  return funnel
}

// This is a blocking call, you'll regularly want to execute it within a goroutine.
func (f *Funnel) Listen() {
  for {
    connection := f.Connections.Get()

    for {
      job, err := connection.Fetch(f.FetchCount, f.FetchTimeout, f.Queues...); if err != nil {
        log.Printf("Error fetching jobs in background: %v\n", err.Error())
        break
      }

      f.Incoming <- job
    }

    connection.Close()
  }
}

// This is a blocking call, you'll regularly want to execute it within a goroutine.
func (f *Funnel) Dispatch() {
  for {
    job := <- f.Outgoing
    connection := f.Connections.Get()
    connection.AddJob(job.Queue, string(job.Payload), "10s") // TODO: Push timeout should be configurable.
    connection.Close()
  }
}

func (f *Funnel) Close() {

}
