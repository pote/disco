package disco

import(
  "testing"
)

func TestConsumerFetch(t *testing.T) {
  pool, _ := NewPool(2, 5, "240s", 1)
  connection := pool.Get()
  connection.AddJob("disco-test-queue", "this-is-the-payload", "10s")

  consumer := pool.NewConsumer("disco-test-queue")
  job, err := consumer.Fetch(1, "10s")

  if err != nil {
    t.Fatal(err)
  }

  if job.ID == "" {
    t.Error("fetched jobs should always have ids")
  }

  if string(job.Payload) != "this-is-the-payload" {
    t.Errorf("Expected payload does not match: '%v'", string(job.Payload))
  }
}

func TestConsumerIncomingChannel(t *testing.T) {
  pool, _ := NewPool(2, 5, "240s", 1)
  connection := pool.Get()
  connection.AddJob("disco-test-queue", "this-is-the-payload", "10s")

  consumer := pool.NewConsumer("disco-test-queue")
  go consumer.FetchIntoChannel(1, "10s")

  job, ok := <- consumer.IncomingJobs

  if !ok {
    t.Fatal("I... I guess something is not ok")
  }

  if job.ID == "" {
    t.Error("fetched jobs should always have ids")
  }

  if string(job.Payload) != "this-is-the-payload" {
    t.Errorf("Expected payload does not match: '%v'", string(job.Payload))
  }
}
