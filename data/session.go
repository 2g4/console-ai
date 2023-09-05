package data

var messages []Message

func AppendMessage(role string, content string) {
	messages = append(messages, Message{Role: role, Content: content})
}

func GetMessages() []Message {
	return messages
}
