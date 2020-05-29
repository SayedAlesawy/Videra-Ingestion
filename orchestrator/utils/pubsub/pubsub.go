package pubsub

// Subscriber An interface for subscriber types
type Subscriber interface {
	Notify(notification interface{}) //A function to notify the subscriber with new events
}
