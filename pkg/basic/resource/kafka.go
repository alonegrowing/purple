package resource

import (
	"errors"
	"fmt"

	"github.com/alonegrowing/purple/pkg/kernel/kafka"
	log "github.com/alonegrowing/purple/pkg/kernel/logging"
)

var (
	KAFKA_CONSUMEFROM_NOT_INIT = errors.New("kafka consume from not init ")
	KAFKA_PRODUCERTO_NOT_INIT  = errors.New("kafka producer to not init or no has live client ")
)

var consumeClientMap map[string]*kafka.KafkaConsumeClient
var producerClientMap map[string]*kafka.KafkaClient
var syncProducerClientMap map[string]*kafka.KafkaSyncClient

func InitKafkaConsume(consumeConfigs []kafka.KafkaConsumeConfig) error {
	if consumeClientMap == nil {
		consumeClientMap = make(map[string]*kafka.KafkaConsumeClient)
	}
	for _, config := range consumeConfigs {
		if _, ok := consumeClientMap[config.ConsumeFrom]; ok {
			continue
		}
		client, err := kafka.NewKafkaConsumeClient(config)
		if err != nil {
			log.Error("rpc.InitKafkaConsume,err:", err)
			return err
		}
		consumeClientMap[config.ConsumeFrom] = client
	}
	return nil
}

func InitKafkaProducer(producerConfigs []kafka.KafkaProductConfig) error {
	producerClientMap = make(map[string]*kafka.KafkaClient)
	syncProducerClientMap = make(map[string]*kafka.KafkaSyncClient)
	for _, config := range producerConfigs {
		if config.UseSync == true {
			fmt.Println(config.ProducerTo, "--")
			if _, ok := syncProducerClientMap[config.ProducerTo]; ok {
				continue
			}

			client, err := kafka.NewSyncProducterClient(config)
			if err != nil {
				log.Error("rpc.InitKafkaProducer,err:", err)
				return err
			}
			syncProducerClientMap[config.ProducerTo] = client
			fmt.Println(config.ProducerTo)
		} else {
			if _, ok := producerClientMap[config.ProducerTo]; ok {
				continue
			}
			client, err := kafka.NewKafkaClient(config)
			if err != nil {
				log.Error("rpc.InitKafkaProducer,err:", err)
				return err
			}
			producerClientMap[config.ProducerTo] = client
		}
	}
	return nil
}

func GetKafkaConsumeClient(consumeFrom string) (*kafka.KafkaConsumeClient, error) {
	if client, ok := consumeClientMap[consumeFrom]; ok {
		return client, nil
	}
	return nil, KAFKA_CONSUMEFROM_NOT_INIT
}

func GetSyncProducerClient(producerTo string) (*kafka.KafkaSyncClient, error) {
	if client, ok := syncProducerClientMap[producerTo]; ok {
		return client, nil
	}
	return nil, KAFKA_PRODUCERTO_NOT_INIT
}

func GetKafkaProducerClient(producerTo string) (*kafka.KafkaClient, error) {
	if client, ok := producerClientMap[producerTo]; ok {
		return client, nil
	}
	return nil, KAFKA_PRODUCERTO_NOT_INIT
}

func CloseAllClient() error {
	for _, client := range consumeClientMap {
		client.Close()
	}
	for _, client := range producerClientMap {
		client.Close()
	}
	for _, client := range syncProducerClientMap {
		client.Close()
	}
	return nil
}
