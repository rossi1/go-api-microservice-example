package internal

//import (
//	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
//)

//type KafkaProducer struct {
//	Producer *kafka.Producer
//	Topic    string
//}

//type KafkaConsumer struct {
//	Consumer *kafka.Consumer
//}

//func NewKafkaProducer(cfg Config) (*KafkaProducer, error) {
//	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.HOST})//
//	if err != nil {
//		return nil, err

//	}

//	return &KafkaProducer{
//		Producer: producer,
//		Topic:    cfg.TOPIC,
//	}, nil
//}

//func (k *KafkaProducer) Close() {
//	k.Producer.Close()
//}

//func NewKaftaConsumer(cfg Config) (*KafkaConsumer, error) {
//	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers": cfg.HOST,
//		"group.id":          "myGroup",
//		"auto.offset.reset": "earliest",
//	})

//	if err != nil {
//		return nil, err
//	}
//	err = consumer.Subscribe(cfg.TOPIC, nil)

//	if err != nil {
//		return nil, err
//	}

//	return &KafkaConsumer{
//		Consumer: consumer,
//	}, nil
//}

//func (k *KafkaConsumer) Close() {
//	k.Consumer.Close()
//}
//*/
