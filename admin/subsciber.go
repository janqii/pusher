package admin

import (
	"encoding/json"
	"fmt"
	"github.com/janqii/pusher/utils"
)

type TopicAndUriInfo struct {
	Topic string
	Uri   string
}

type McpackKeyCopyInfo struct {
	From string
	To   string
}

type UbrpcInfo struct {
	ServiceName string
	Method      string
	CmdKey      string
}

type ReqCheckInfo struct {
	PassWhenNoCheckSegment int
	Expression             string
}

type LocalAddr struct {
	Addr string
	Tag  string
}

type WebfootInfo struct {
	Name string
	Tag  string
}

type MachineInfo struct {
	Local  []LocalAddr
	Naming WebfootInfo
}

type Subscriber struct {
	ConsumerType     int // 0:竞争   1:多主
	ConsumerProt     int // 0:http   1:nshead
	ConsumerConveter int // 0:json   1:mcpack1   2:mapack2   3:msgpack
	PushRetryTimeMs  int
	RetryTimes       int
	PushDelayTimeMs  int
	ReqCheck         ReqCheckInfo
	TopicAndUri      []TopicAndUriInfo
	McpackKeyCopy    []McpackKeyCopyInfo
	Ubrpc            UbrpcInfo
	Machines         MachineInfo
}

type SubscribeManager struct {
	ZkClient      *utils.ZK
	SubscriberMap map[string]*Subscriber
	SubscriberNum int
}

func (m *SubscribeManager) Startup() error {
	fmt.Println("SubscribeManager startup")
	return nil
}

func (m *SubscribeManager) Shutdown() {
	fmt.Println("SubscribeManager shutdown")
}
