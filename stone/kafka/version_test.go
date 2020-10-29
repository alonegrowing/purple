package kafka

import (
	"log"
	"testing"

	"github.com/Shopify/sarama"
)

func TestCheckKafakHeadersSupported(t *testing.T) {
	// 0.11 cluster
	hosts := []string{"192.168.40.22:9092", "ali-a-inf-kafka-test11.bj:9092", "ali-c-inf-kafka-test12.bj:9092", "ali-a-inf-kafka-test13.bj:9092"}
	supported, err := checkKafkaHeadersSupported(hosts)
	log.Printf("checkKafkaHeadersSupported supported=%t, error=%v", supported, err)

	// 0.10 cluster
	hosts = []string{"192.168.40.22:9092", "10.10.0.102:9092", "10.10.0.103:9092"}
	supported, err = checkKafkaHeadersSupported(hosts)
	log.Printf("checkKafkaHeadersSupported supported=%t, error=%v", supported, err)
}

func TestAdjustProducerVersion(t *testing.T) {
	hosts := []string{"ali-a-inf-kafka-test11.bj:9092", "ali-c-inf-kafka-test12.bj:9092", "ali-a-inf-kafka-test13.bj:9092"}
	conf := sarama.NewConfig()
	adjustProducerVersion(hosts, conf)
}

func TestAdjustConsumerVersion(t *testing.T) {
	conf := sarama.NewConfig()
	adjustConsumerVersion("ali-a-inf-zk04.bj:2181,ali-a-inf-zk05.bj:2181,ali-a-inf-zk06.bj:2181,ali-a-inf-zk07.bj:2181/config/inke/inf/mq/kafka02", conf)
}
