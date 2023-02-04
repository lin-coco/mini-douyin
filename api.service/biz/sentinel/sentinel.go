package sentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func InitSentinel() {
	err := sentinel.InitWithConfigFile("./sentinel.yaml")
	if err != nil {
		hlog.Fatalf("Unexpected error: %+v", err)
	}

	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "api_service",
			Threshold:              50,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
		{
			Resource:               "core_service",
			Threshold:              3,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		}, {
			Resource:               "interaction_service",
			Threshold:              3,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		}, {
			Resource:               "society_service",
			Threshold:              3,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	if err != nil {
		hlog.Fatalf("Unexpected error: %+v", err)
		return
	}
}
