package gorabbitmq_starter

import (
	goframeworkrabbitmq "github.com/kordar/goframework-rabbitmq"
	logger "github.com/kordar/gologger"
	"github.com/spf13/cast"
)

type RabbitmqPublishModule struct {
	name string
	load func(moduleName string, itemId string, item map[string]string)
}

func NewRabbitmqPublishModule(name string, load func(moduleName string, itemId string, item map[string]string)) *RabbitmqPublishModule {
	return &RabbitmqPublishModule{name, load}
}

func (m RabbitmqPublishModule) Name() string {
	return m.name
}

func (m RabbitmqPublishModule) _load(id string, cfg map[string]string) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	err := goframeworkrabbitmq.AddPublishInstance(id, cfg["dsn"])
	if err != nil {
		logger.Errorf("[gorabbitmq-publish-starter] 初始化rabbitmq异常，err=%v", err)
		return
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		logger.Debugf("[%s] triggering custom loader completion", m.Name())
	}

	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
}

func (m RabbitmqPublishModule) Load(value interface{}) {

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

func (m RabbitmqPublishModule) Close() {
}
