package subscriber

// Subscriber is the interface for configcenter subscribers.
type Subscriber interface {
	// AddListener adds a listener to the subscriber.
	AddListener(listener func(key, data string)) error
	// Value returns the value of the subscriber.
	Value() (string, error)
	Key() string
}