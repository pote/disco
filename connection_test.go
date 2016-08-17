package disco

import(
  "testing"
  "time"

  "github.com/garyburd/redigo/redis"
)

// This test will attempt a connection to a Disque server specified in
// the DISQUE_NODES environmnet variable.
func TestNewConnectionToAvailableNode(t *testing.T) {
  connection, err := NewConnection(1)

  if err != nil {
    t.Fatal(err)
  }

  response, err := redis.String(connection.Do("PING"))

  if err != nil {
    t.Fatal(err)
  }

  if response != "PONG" {
    t.Error("Expected PONG response from Disque")
  }
}

func TestAddJob(t *testing.T) {
  connection, _ := NewConnection(1)

  id, err := connection.AddJob("disco-test-queue", "this-is-the-payload", "10s")

  if err != nil {
    t.Fatal(err)
  }

  if id == "" {
    t.Error("No job id returned from ADDJOB")
  }
}

func TestConnectionFetch(t *testing.T) {
  connection, err := NewConnection(1)
  connection.AddJob("disco-test-queue", "this-is-the-payload", "10s")

  job, err := connection.Fetch(1, time.Second * 10, "disco-test-queue")

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
