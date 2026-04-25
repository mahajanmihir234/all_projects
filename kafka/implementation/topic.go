package implementation

//
// Topic -> id, messages[], generate_log(), get_next_message()
// Consumer Group -> id, consume(topic_message) -> void, read_topic_again(topic)
// KafkaSystem -> add_topic, remove_topic, add_consumer_to_topic, remove_consumer_to_topic

type TopicMessage struct {
	message string
}

type Topic struct {
	id       string
	messages []TopicMessage
}

func (t Topic) Id() string { return t.id }

func (t Topic) GetAllMessages() []TopicMessage {
	return t.messages
}
