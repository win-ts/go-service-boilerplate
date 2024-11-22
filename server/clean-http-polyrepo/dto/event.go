package dto

// Event represents the event message to be produced
type Event struct {
	Event   string       `json:"event"`
	Payload EventPayload `json:"payload"`
}

// EventPayload represents the payload of the event message
type EventPayload struct {
	Data string `json:"data"`
}
