package disco

import(
  "log"
  "time"
)

// A funnel is a high-level API for Disque usage: it acts as a bridge between Disque and
// Go native channels, allowing for idiomatic interaction with the datastore.
type Funnel struct {
  Queues []string
  Incoming chan Job
  Outgoing chan Job
  Connections *Pool

  FetchTimeout time.Duration
  FetchCount int

  Closed bool
}

// Creates a new funnel with a specific queue configuration and starts the
// appropriate goroutines to keep it's go channels synchronized with Disque.
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

// Takes a connection from the funnel's Disque connection pool and uses it to fetch
// jobs from the funnel's configured queues.
//
// This is a blocking call which is launched on a goroutine when #NewFunnel is called,
// you won't reguarly call it directly, but it's left as a public method to allow
// more flexibility of use cases.
func (f *Funnel) Listen() {
  for {
    if f.Closed {
        close(f.Incoming)
        return
    }
    connection := f.Connections.Get()

    for {
      job, err := connection.GetJob(f.FetchCount, f.FetchTimeout, f.Queues...); if err != nil {
        log.Printf("Error fetching jobs in background: %v\n", err.Error())
        break
      }

      if f.Closed {
        connection.NAck(job.ID)
        connection.Close()
        close(f.Incoming)
        return
      }

      f.Incoming <- job
      connection.Ack(job.ID)
    }

    connection.Close()
  }
}

// Listens to the `Outgoing` channel in the funnel, and dispatches any messages received
// to it's appropriate queue taking a connection from the funnel's pool.
//
// This is a blocking call which is launched on a goroutine when #NewFunnel is called,
// you won't reguarly call it directly, but it's left as a public method to allow
// more flexibility of use cases.
//
func (f *Funnel) Dispatch() {
  for {
    select {
    case job := <- f.Outgoing:
      connection := f.Connections.Get()
      connection.AddJob(job.Queue, string(job.Payload), time.Second * 10) // TODO: Push timeout should be configurable.
      connection.Close()

      if f.Closed {
        close(f.Outgoing)
        return
      }
    case <- time.Tick(time.Second):
      if f.Closed {
        close(f.Outgoing)
        return
      }
    }
  }
}

// Marks the funnel as closed, which in turn closes its internal go channels
// gracefully.
func (f *Funnel) Close() {
  f.Closed = true
}
