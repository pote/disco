package disco

// Wraps a Disque job in a Go struct.
type Job struct {
  Queue   string
  ID      string
  Payload []byte
}
