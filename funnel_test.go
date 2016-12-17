package disco

import(
  "testing"
  "time"
)


func TestFunnelIncomingChannel(t *testing.T) {
  pool, _ := NewPool(2, 5, 1, time.Second * 200)
  connection := pool.Get()
  connection.AddJob("disco-test-queue", "this-is-the-payload", time.Second * 10)
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

func TestFunnelCloseBehaviour(t *testing.T) {
  pool, _ := NewPool(2, 5, 1, time.Second * 200)
  funnel := pool.NewFunnel("disco-test-queue")

  defer func() {
    if r := recover(); r != nil {
      t.Error("We shouldn't be panicking under any circumstances.")
    }
  }()

  funnel.Close()

  select {
  case funnel.Outgoing <- Job{Queue: "disco-test-queue", Payload: []byte("this-is-the-payload")}:
    // NoOp, this should work.
  case <- time.Tick(time.Second):
    t.Error("we shouldnt be blocking the outgoing channel immediately after closing")
  }

  select {
  case <- funnel.Incoming:
    // NoOp, we want to fail gracefully at this point.
  case <- time.Tick(time.Second):
    t.Error("Closed funnels blocking external sends")
  }

  time.Sleep(time.Second * 2)
  select {
  case <- funnel.Incoming:
    // NoOp, we want to fail gracefully at this point.
  case <- time.Tick(time.Second):
    t.Error("Closed funnels blocking reads")
  }

  select {
  case funnel.Outgoing <- Job{Queue: "disco-test-queue", Payload: []byte("this-is-the-payload")}:
    // NoOp, we want to fail gracefully at this point.
  case <- time.Tick(time.Second):
    t.Error("Closed funnels blocking sends")
  }
}
