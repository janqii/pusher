package global

import (
	"github.com/janqii/pusher/admin"
	"gitlab.baidu.com/go/sarama"
)

var (
	KafkaClient *sarama.Client
	SubManager  *admin.SubscriberManager
)
