package main

import (
	"errors"
	"fmt"
	"sync"
)

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

type KafkaSystem struct {
	topics             map[*Topic][]Consumer
	mutex              sync.Mutex
	consumerToLocksMap map[Consumer]*sync.Mutex
}

func (s *KafkaSystem) AddTopic(topic Topic) error {
	if _, ok := s.topics[&topic]; ok {
		return errors.New("topic already exists")
	}

	s.topics[&topic] = []Consumer{}
	return nil
}

func (s *KafkaSystem) RemoveTopic(topic Topic) error {
	if _, ok := s.topics[&topic]; ok {
		delete(s.topics, &topic)
		return nil
	}
	return errors.New("Topic not found")
}

func (s *KafkaSystem) AddConsumersToTopic(topic Topic, consumers []Consumer) error {
	if _, ok := s.topics[&topic]; ok {
		s.topics[&topic] = append(s.topics[&topic], consumers...)
		return nil
	}
	return errors.New("Topic not found")
}

func (s *KafkaSystem) RemoveConsumersFromTopic(topic Topic, consumers []Consumer) error {
	if _, ok := s.topics[&topic]; ok {
		// newConsumers :=
		// for _, consumer := range s.topics[&topic] {

		// }
		// s.topics[&topic] = append(s.topics[&topic], consumers...)
		return nil
	}
	return errors.New("Topic not found")
}

func (s *KafkaSystem) AddMessageToTopic(topic Topic, message TopicMessage) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.topics[&topic]; !ok {
		return errors.New("Topic not found")
	}
	topic.messages = append(topic.messages, message)

	for _, consumer := range s.topics[&topic] {
		mutex := s.consumerToLocksMap[consumer]
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			consumer.Consume(message)
		}()
	}
	return nil
}

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

func main() {}
