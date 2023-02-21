package vd

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
)

func InitVD() {

	binding.SetLooseZeroMode(true)
	//数值不为负数
	binding.MustRegValidateFunc("NotNegative", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s := args[0].(float64)
		if s < 0 {
			return fmt.Errorf("the args can not be less 0")
		}
		return nil
	})
	//id >= 0
	binding.MustRegValidateFunc("GreaterZero", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s := args[0].(float64)
		if s <= 0 {
			return fmt.Errorf("the args can not be less or equal 0")
		}
		return nil
	})

	//字符串不为默认值
	binding.MustRegValidateFunc("NotStringDefault", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s := args[0].(string)
		if s == "" {
			return fmt.Errorf("the args can not be \"\"")
		}
		return nil
	})

	//用户名密码标题长度不能超过20
	binding.MustRegValidateFunc("LessTextNotNull", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s := args[0].(string)
		if s == "" {
			return fmt.Errorf("the args can not be \"\"")
		}
		if len(s) > 20 {
			return fmt.Errorf("the args too large")
		}
		return nil
	})
	//评论消息长度不能超过40
	binding.MustRegValidateFunc("LargeTextNotNull", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s := args[0].(string)
		if s == "" {
			return fmt.Errorf("the args can not be \"\"")
		}
		if len(s) > 40 {
			return fmt.Errorf("the args too large")
		}
		return nil
	})

	//用户名密码标题长度不能超过20
	binding.MustRegValidateFunc("LessText", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s := args[0].(string)
		if s == "" {
			return nil
		}
		if len(s) > 20 {
			return fmt.Errorf("the args too large")
		}
		return nil
	})
	//评论消息长度不能超过40
	binding.MustRegValidateFunc("LargeText", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s := args[0].(string)
		if s == "" {
			return nil
		}
		if len(s) > 40 {
			return fmt.Errorf("the args too large")
		}
		return nil
	})

	binding.MustRegValidateFunc("NotNil", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		bytes := args[0].([]byte)
		if len(bytes) == 0 {
			return fmt.Errorf("the args can not be nil")
		}
		return nil
	})
}
