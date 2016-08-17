package disco

import(
  "log"
)

type Funnel struct {
  Queues []string
  Incoming chan Job
  Outgoing chan Job
  Connections *Pool
}

func NewFunnel(pool *Pool, queues ...string) Funnel {
  incoming := make(chan Job)
  outgoing := make(chan Job)

  funnel := Funnel{Queues: queues, Incoming: incoming, Outgoing: outgoing, Connections: pool}

  go funnel.Dispatch()
  return funnel
}

// This is a blocking call, you'll regularly want to execute it within a goroutine.
func (f *Funnel) Listen(count int, fetchTimeout string) {
  for {
    connection := f.Connections.Get()

    for {
      job, err := connection.Fetch(count, fetchTimeout, f.Queues...); if err != nil {
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
