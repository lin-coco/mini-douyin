package log

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
)

func InitLog() {
	file, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	hlog.SetOutput(file)
}
