// Disco is a flexible, idiomatic approach to a Go Disque client.
//
// It attempts to provide two ways of using Disque: a low level API that allows more flexibility and control for users that need it and high level API for the most common usage patterns.
//
//
// High Level Api
//
// The high level API attempts to provide a common usage pattern in a idiomatic Go manner, ideally it should simplify Disque usage by not having to deal with the nuts and bolts of the low level API.
//
// Funnels
//
// Funnels are an abstraction on top of a `disco.Pool`: they provide Go channels that you can use to enqueue or receive jobs from Disque.
//
//
//   // See GoDoc for further details in connection Pool creation.
//   pool, _ := NewPool(2, 5, 1, time.Second * 200)
//   funnel := pool.NewFunnel("disco-test-queue", "other-queue")
//   defer funnel.Close()
//
//   // Enqueue jobs simply by directing them to the Outgoing channel.
//   funnel.Outgoing <- Job{Queue: "disco-test-queue", Payload: []byte("this-is-the-payload")}:
//
//   // Receive jobs from disque simply by leveraging the Incoming channel, you can leverage
//   // common Go patterns such as a select statement to handle timeouts or other kinds of errors.
//   select {
//   case job, ok := <- funnel.Incoming:
//     string(job.Payload) //=> "this-is-the-payload" {
//   case <- time.Tick(time.Second):
//     // Handle timeout (or not)
//   }
//
// A funnel will also manage the job's lifecycle for you: jobs received via the `Incoming` channel will be acknowledged in Disque automatically (you'll still have the option to put it back in the queue if need be) and jobs fetched from Disque after the funnel is closed will be automatically NAcked so as not to lose data.
//
package disco // import "github.com/pote/disco"
