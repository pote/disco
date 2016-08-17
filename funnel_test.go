package disco

import(
  "testing"
  "time"
)


func TestFunnelIncomingChannel(t *testing.T) {
  pool, _ := NewPool(2, 5, 1, time.Second * 200)
  connection := pool.Get()
  connection.AddJob("disco-test-queue", "this-is-the-payload", "10s")
  connection.Close()

  funnel := pool.NewFunnel("disco-test-queue")

  select {
  case job, ok := <- funnel.Incoming:
    if !ok {
      t.Fatal("I... I guess something is not ok")
    }

    if job.ID == "" {
      t.Error("fetched jobs should always have ids")
    }

    if string(job.Payload) != "this-is-the-payload" {
      t.Errorf("Expected payload does not match: '%v'", string(job.Payload))
    }
  case <- time.Tick(time.Second):
    t.Error("Failed to fetch job in a timely manner")
  }
}

func TestFunnelOutgoingChannel(t *testing.T) {
  pool, _ := NewPool(2, 5, 1, time.Second * 200)

  funnel := pool.NewFunnel("disco-test-queue")
  funnel.Outgoing <- Job{Queue: "disco-test-queue", Payload: []byte("this-is-the-payload")}

  select {
  case job, ok := <- funnel.Incoming:
    if !ok {
      t.Fatal("I... I guess something is not ok")
    }

    if job.ID == "" {
      t.Error("fetched jobs should always have ids")
    }

    if string(job.Payload) != "this-is-the-payload" {
      t.Errorf("Expected payload does not match: '%v'", string(job.Payload))
    }
  case <- time.Tick(time.Second):
    t.Error("Failed to fetch job in a timely manner")
  }
}
