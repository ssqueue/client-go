package queue

type InputMessage struct {
	Data       string `json:"data,omitempty"`
	Persistent bool   `json:"persistent,omitempty"`
}

type OutputMessage struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}
