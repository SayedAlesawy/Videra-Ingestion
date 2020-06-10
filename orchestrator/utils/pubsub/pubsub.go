package pubsub

// Subscriber An interface for subscriber types
type Subscriber interface {
	Notify(event Event) //A function to notify the subscriber with new events
}

// Event Represents an event in the pub-sub notifcation
type Event struct {
	Type string        //Type of the broadcasted event
	Args []interface{} //Optional event args
}

// NewEvent A function to create a new event
func NewEvent(eventType string, args ...interface{}) Event {
	return Event{
		Type: eventType,
		Args: args,
	}
}
