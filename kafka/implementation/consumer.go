package implementation

import "fmt"

type Consumer interface {
	Consume(topicMessage TopicMessage) error
	ConsumeTopic(topic Topic) error
}

type PrintConsumer struct {
	consumedMessages []TopicMessage
}

func (p PrintConsumer) Consume(topicMessage TopicMessage) error {
	fmt.Println(topicMessage.message)
	p.consumedMessages = append(p.consumedMessages, topicMessage)
	return nil
}

func (p PrintConsumer) ConsumeTopic(topic Topic) error {
	p.consumedMessages = []TopicMessage{}
	for _, message := range topic.messages {
		err := p.Consume(message)
		if err != nil {
			return err
		}
	}
	return nil
}
