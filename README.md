# Disco - [![GoDoc](https://godoc.org/github.com/pote/disco?status.svg)](https://godoc.org/github.com/pote/disco)

A flexible, idiomatic approach to a Go [Disque](https://github.com/antirez/disque) client.

### The Project

Disco attempts to provide two ways of using Disque: a [low level API](#low-level-api) that allows more flexibility and control for users that need it and [high level API](#high-level-api) for the most common usage patterns.

## High Level Api

The high level API attempts to provide a common usage pattern in a idiomatic Go manner, ideally it should simplify Disque usage by not having to deal with the nuts and bolts of the low level API.

### Funnels

Funnels are an abstraction on top of a [`disco.Pool`](https://godoc.org/github.com/pote/disco#Pool): they provide Go channels that you can use to enqueue or receive jobs from Disque.


```Go
  // See GoDoc for further details in connection Pool options.
  pool, _ := NewPool(2, 5, 1, time.Second * 200)
  funnel := pool.NewFunnel("disco-test-queue", "other-queue")
  defer funnel.Close()

  // Enqueue jobs simply by directing them to the Outgoing channel.
  funnel.Outgoing <- Job{Queue: "disco-test-queue", Payload: []byte("this-is-the-payload")}:

  // Receive jobs from disque simply by leveraging the Incoming channel, you can leverage
  // common Go patterns such as a select statement to handle timeouts or other kinds of errors.
  select {
  case job, ok := <- funnel.Incoming:
    string(job.Payload) //=> "this-is-the-payload" {
  case <- time.Tick(time.Second):
    // Handle timeout (or not)
  }
```

A funnel will also manage the job's lifecycle for you: jobs received via the `Incoming` channel will be acknowledged in Disque automatically (you'll still have the option to put it back in the queue if need be) and jobs fetched from Disque after the funnel is closed will be automatically NAcked so as not to lose data.

## Low Level Api

### Connections

Connections represent a persistent connection to a Disque cluster, it's the most basic form of Disco usage there is. Disco is built on top of [Redigo](https://github.com/garyburd/redigo), and the Connection struct is a Disque-specific wrapper around a [redigo Conn interface](https://godoc.org/github.com/garyburd/redigo/redis#Conn), which means you can send commands to Disque directly.

```Go
  // Will connect to the Disque nodes specified in the $DISQUE_NODES env variable.
  connection, err := disco.NewConnection(100)

  connection, err := disco.NewConnectionToURLS(100, "localhost:7701,localhost:7702,localhost:7703")

  connection.Do("PING")
```

### Connection Pools

In most cases it's better to have a global connection pool that your application uses instead of manually creating them each time.

```Go
  // Will connect to the Disque nodes specified in the $DISQUE_NODES env variable.
  // Args are: Max idle connections, max active, cycle and idle timeout.
  // see GoDoc for further details
  pool, err := disco.NewPool(2, 5, 1, time.Second * 200)

  connection := pool.Get()
  defer connection.Close()

  connection.Do("PING")
```

### Wrappers for Disque commands.

#### AddJob

The [`AADDJOB`](https://github.com/antirez/disque#addjob-queue_name-job-ms-timeout-replicate-count-delay-sec-retry-sec-ttl-sec-maxlen-count-async) command is one of the two most used one, it enqueues a payload in a given queue in Disque.

```Go
  connection, _ := disco.NewConnection(100)

  id, err := connection.AddJob("disco-test-queue", "this-is-the-payload", time.Second * 10)
```

#### GetJob

[`GETJOB`](https://github.com/antirez/disque#getjob-nohang-timeout-ms-timeout-count-count-withcounters-from-queue1-queue2--queuen) is the other fundamental Disque command: fetches enqueued jobs from a list of specified queues.

Keep in mind that this is a blocking call.

```Go
  id, _ := connection.AddJob("disco-test-queue", "this-is-the-payload", time.Second * 10)
  job, err := connection.GetJob(1, time.Second * 10, "disco-test-queue")

  string(job.Payload) //=> "this-is-the-payload"
```

#### Ack

Wrapper around the ['ACKJOB'](https://github.com/antirez/disque#ackjob-jobid1-jobid2--jobidn) command.

Acknowledges the execution of one or more jobs via job IDs


```Go
  job, _ := connection.GetJob(1, time.Second * 10, "disco-test-queue")

  connection.Ack(job.ID)
```



#### NAck

Wrapper around the ['NACK'](https://github.com/antirez/disque#nack-job-id--job-id) command.

The NACK command tells Disque to put the job back in the queue ASAP

```Go
  job, _ := connection.GetJob(1, time.Second * 10, "disco-test-queue")

  connection.NAck(job.ID)
```

## Contributing

You'll need [gpm](https://github.com/pote/gpm) installed in order to pull in the necessary dependencies.

```bash
$ git clone git@github.com:pote/disco.git && cd disco
$ source .env.sample # feel free to cp it to .env and make any config changes you deem necessary.

$ make # Will pull dependencies if necessary, build the project and run the test suite.
```
