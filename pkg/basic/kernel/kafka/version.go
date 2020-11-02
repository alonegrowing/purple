package kafka

import (
	"errors"
	"fmt"
	"time"

	log "github.com/alonegrowing/purple/pkg/basic/kernel/logging"

	"github.com/Shopify/sarama"
	kazoo "github.com/wvanbergen/kazoo-go"
)

const (
	versionRequestReadTimeout  = 10 * time.Second
	versionRequestWriteTimeout = 3 * time.Second

	producerHeaderSupportedVersion = 3
	fetchHeaderSupportedVersion    = 4

	produceRequestKey = 0
	fetchRequestKey   = 1
)

var errConncetedTimeout = errors.New("adjust broker version, connect broker timedout")

func checkContinue(err error, index int, addresses []string) (bool, error) {
	addressesLen := len(addresses)
	if err != nil {
		if index < addressesLen-1 {
			return true, nil
		}
		errMessage := fmt.Errorf("Adjust kafka server error, is Your broker reachable, brokers %v, err %s", addresses, err)
		log.GenLog(errMessage.Error())
		return false, errMessage
	}
	return false, nil
}

// checkKafkaHeadersSupported https://www.confluent.io/blog/upgrading-apache-kafka-clients-just-got-easier/
// https://cwiki.apache.org/confluence/display/KAFKA/KIP-35+-+Retrieving+protocol+version#KIP-35-Retrievingprotocolversion-Aclientdeveloperwantstoaddsupportforanewfeature
// API version request, this requires a >= 0.10.0.0 broker and will cause a disconnect on brokers 0.8.x . NOTE: Due to a bug in broker version 0.9.0.0 & 0.9.0.1 the broker will not close the connection when receiving the API version request, instead the request will time out
func checkKafkaHeadersSupported(addresses []string) (bool, error) {
	for idx, address := range addresses {
		brokers := sarama.NewBroker(address)
		conf := sarama.NewConfig()

		supported := 0
		var resp *sarama.ApiVersionsResponse

		conf.Net.DialTimeout = versionRequestWriteTimeout
		conf.Net.ReadTimeout = versionRequestReadTimeout
		conf.Net.WriteTimeout = versionRequestWriteTimeout
		conf.Version = sarama.V0_11_0_0
		err := brokers.Open(conf)
		ok, err := checkContinue(err, idx, addresses)
		if err != nil {
			return false, err
		}
		if ok {
			continue
		}

		startTime := time.Now()
		var connected bool
		for {
			connected, err = brokers.Connected()
			if err != nil || time.Since(startTime) > versionRequestWriteTimeout {
				if err == nil {
					err = errConncetedTimeout
				}
				goto TRY_NEXT
			}
			if connected {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		resp, err = brokers.ApiVersions(new(sarama.ApiVersionsRequest))
		brokers.Close()
		if err != nil {
			return false, err
		}
		for _, version := range resp.ApiVersions {
			if version.ApiKey == fetchRequestKey && version.MaxVersion >= fetchHeaderSupportedVersion {
				supported |= 1
				continue
			}
			if version.ApiKey == produceRequestKey && version.MaxVersion >= producerHeaderSupportedVersion {
				supported |= 1 << 1
				continue
			}
		}
		return supported == 3, nil
	TRY_NEXT:
		brokers.Close()
		ok, err = checkContinue(err, idx, addresses)
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func getKafkaBrokerListFromZK(zkHosts string) ([]string, error) {
	var zookNodes []string
	zooConfig := kazoo.NewConfig()
	zookNodes, zooConfig.Chroot = kazoo.ParseConnectionString(zkHosts)
	var kz *kazoo.Kazoo
	var err error
	if kz, err = kazoo.NewKazoo(zookNodes, zooConfig); err != nil {
		return nil, err
	}
	brokers, err := kz.BrokerList()
	kz.Close()
	if err != nil {
		return nil, err
	}
	return brokers, nil
}

func adjustProducerVersion(brokers []string, conf *sarama.Config) error {
	supported, err := checkKafkaHeadersSupported(brokers)
	if err != nil {
		log.GenLogf("adjustProducerVersion brokers:%v, error %s", brokers, err)
		return err
	}
	log.GenLogf("adjustProducerVersion brokers:%v, supportedv0_11:%t", brokers, supported)
	log.Infof("adjustProducerVersion brokers:%v, supportedv0_11:%t", brokers, supported)
	if supported {
		conf.Version = sarama.V0_11_0_0
	}
	return nil
}

func adjustConsumerVersion(zkHosts string, conf *sarama.Config) error {
	brokers, err := getKafkaBrokerListFromZK(zkHosts)
	if err != nil {
		log.GenLogf("adjustConsumerVersion get kafka brokers error (zkhosts:%s), error %s", zkHosts, err)
		return err
	}
	supported, err := checkKafkaHeadersSupported(brokers)
	if err != nil {
		return err
	}
	log.GenLogf("adjustConsumerVersion zkhosts:%s, supportedv0_11:%t", zkHosts, supported)
	log.Infof("adjustConsumerVersion zkhosts:%s, supportedv0_11:%t", zkHosts, supported)
	if supported {
		conf.Version = sarama.V0_11_0_0
	}
	return nil
}
