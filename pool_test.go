package disco

import(
  "testing"

  "github.com/garyburd/redigo/redis"
)

// This test will attempt a connection pool to a Disque server specified in
// the DISQUE_NODES environmnet variable.
func TestNewPooltoAvailableNode(t *testing.T) {
  pool, err := NewPool(2, 2, "240s", 1)

  if err != nil {
    t.Fatal(err)
  }

  connection := pool.Get()
  defer connection.Close()

  response, err := redis.String(connection.Do("PING"))

  if err != nil {
    t.Fatal(err)
  }

  if response != "PONG" {
    t.Error("Expected PONG response from Disque")
  }
}
