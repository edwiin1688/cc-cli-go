package api

import (
	"context"
	"encoding/json"
)

type MockClient struct {
	Events []Event
	Error  error
}

type Event struct {
	Type    string
	Content string
	Delta   string
}

func NewMockClient() *MockClient {
	return &MockClient{
		Events: []Event{},
	}
}

func (m *MockClient) Stream(ctx context.Context, req *Request) (<-chan StreamEvent, error) {
	ch := make(chan StreamEvent)

	if m.Error != nil {
		close(ch)
		return ch, m.Error
	}

	go func() {
		defer close(ch)

		for _, event := range m.Events {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- StreamEvent{
					Type:  event.Type,
					Delta: json.RawMessage([]byte(`{"type":"` + event.Delta + `"}`)),
				}
			}
		}

		ch <- StreamEvent{Type: "message_stop"}
	}()

	return ch, nil
}

func (m *MockClient) AddTextEvent(text string) {
	m.Events = append(m.Events, Event{
		Type:    "content_block_delta",
		Content: text,
		Delta:   "text_delta",
	})
}

func (m *MockClient) SetError(err error) {
	m.Error = err
}
