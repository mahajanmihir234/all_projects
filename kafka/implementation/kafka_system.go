package implementation

import (
	"errors"
	"sync"
)

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

// func (s *KafkaSystem) ConsumeMessages(topic Topic) error {
// 	if _, ok := s.topics[&topic]; !ok {
// 		return errors.New("Topic not found")
// 	}

// 	for _, message := range topic.messages {

// 	}

// 	return nil
// }
