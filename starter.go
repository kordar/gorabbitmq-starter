package gorabbitmq_starter

import (
	goframeworkrabbitmq "github.com/kordar/goframework-rabbitmq"
	logger "github.com/kordar/gologger"
	"github.com/spf13/cast"
)

type RabbitmqModule struct {
	name string
	args map[string]interface{}
	load func(moduleName string, itemId string, item map[string]string)
}

func NewRabbitmqModule(name string, load func(moduleName string, itemId string, item map[string]string), args map[string]interface{}) *RabbitmqModule {
	return &RabbitmqModule{name, args, load}
}

func (m RabbitmqModule) Name() string {
	return m.name
}

func (m RabbitmqModule) _load(id string, cfg map[string]string) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	err := goframeworkrabbitmq.AddRabbitmqInstanceArgs(id, cfg, m.args)
	if err != nil {
		logger.Errorf("[gorabbitmq-starter] 初始化rabbitmq异常，err=%v", err)
		return
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		logger.Debugf("[%s] triggering custom loader completion", m.Name())
	}

	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
}

func (m RabbitmqModule) Load(value interface{}) {
	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, cast.ToStringMapString(value))
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMapString(item))
	}

}

func (m RabbitmqModule) Close() {
}
