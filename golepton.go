package golepton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const modelURL = "https://jetmoe-8b-chat.lepton.run/api/v1/chat/completions"
const defaultModel = "jetmoe/jetmoe-8b-chat"

type Model struct {
	APIToken string
	URL      string
	payload  Payload
}

func New(token string) *Model {
	return &Model{
		APIToken: token,
		URL:      modelURL,
		payload: Payload{
			Model:       defaultModel,
			Messages:    make([]Message, 0),
			Temperature: 0.7,
		},
	}
}

func (m *Model) Complete(prompt string) (string, error) {
	res, err := m.req(prompt)
	if err != nil {
		return "", err
	}
	ans := res.Choices[rand.Intn(len(res.Choices))]
	return ans.Message.Content, nil
}

func (m *Model) req(content string) (completion, error) {
	m.payload.Messages = append(m.payload.Messages, mkMessage(content))
	data := m.payload
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return completion{}, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", m.URL, body)
	if err != nil {
		return completion{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.APIToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return completion{}, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return completion{}, err
	}
	var c completion
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return completion{}, err
	}

	return c, nil
}
