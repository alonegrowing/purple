package resource

import (
	"errors"
	"fmt"

	"git.inke.cn/BackendPlatform/golang/kafka"
	log "git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/BackendPlatform/golang/utils"
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

		st := utils.NewServiceStatEntry(kafka.C_KAFKA_PRE, config.ConsumeFrom)

		client, err := kafka.NewKafkaConsumeClient(config)

		if err != nil {
			log.Error("rpc.InitKafkaConsume,err:", err)
			st.End(kafka.KAFKA_INIT, kafka.KafkaConsumeInitError)
			return err
		}
		consumeClientMap[config.ConsumeFrom] = client
		st.End(kafka.KAFKA_INIT, kafka.KafkaSuccess)
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
			st := utils.NewServiceStatEntry(kafka.P_KAFKA_PRE, config.ProducerTo)

			client, err := kafka.NewSyncProducterClient(config)
			if err != nil {
				log.Error("rpc.InitKafkaProducer,err:", err)
				st.End(kafka.KAFKA_INIT, kafka.KafkaProducerInitError)
				return err
			}
			syncProducerClientMap[config.ProducerTo] = client
			fmt.Println(config.ProducerTo)
			st.End(kafka.KAFKA_INIT, kafka.KafkaSuccess)

		} else {
			if _, ok := producerClientMap[config.ProducerTo]; ok {
				continue
			}
			st := utils.NewServiceStatEntry(kafka.P_KAFKA_PRE, config.ProducerTo)

			client, err := kafka.NewKafkaClient(config)
			if err != nil {
				log.Error("rpc.InitKafkaProducer,err:", err)
				st.End(kafka.KAFKA_INIT, kafka.KafkaProducerInitError)
				return err
			}
			producerClientMap[config.ProducerTo] = client
			st.End(kafka.KAFKA_INIT, kafka.KafkaSuccess)
		}
	}
	return nil
}

func GetKafkaConsumeClient(consumeFrom string) (*kafka.KafkaConsumeClient, error) {

	stCode := 0
	st := utils.NewServiceStatEntry(kafka.C_KAFKA_PRE, consumeFrom)

	defer func() {
		st.End(kafka.KAFKA_GET_CONSUME_CLIENT, stCode)
	}()

	if client, ok := consumeClientMap[consumeFrom]; ok {
		return client, nil
	}

	stCode = kafka.KafkaGetConsumeClientError
	return nil, KAFKA_CONSUMEFROM_NOT_INIT
}

func GetSyncProducerClient(producerTo string) (*kafka.KafkaSyncClient, error) {

	fmt.Println(producerTo)

	stCode := 0
	st := utils.NewServiceStatEntry(kafka.P_KAFKA_PRE, producerTo)

	defer func() {
		st.End(kafka.KAFKA_GET_PRODUCER_CLIENT, stCode)
	}()

	if client, ok := syncProducerClientMap[producerTo]; ok {
		return client, nil
	}

	stCode = kafka.KafkaGetProducerClientError

	return nil, KAFKA_PRODUCERTO_NOT_INIT
}

func GetKafkaProducerClient(producerTo string) (*kafka.KafkaClient, error) {

	stCode := 0
	st := utils.NewServiceStatEntry(kafka.P_KAFKA_PRE, producerTo)

	defer func() {
		st.End(kafka.KAFKA_GET_PRODUCER_CLIENT, stCode)
	}()

	if client, ok := producerClientMap[producerTo]; ok {
		return client, nil
	}

	stCode = kafka.KafkaGetProducerClientError

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
