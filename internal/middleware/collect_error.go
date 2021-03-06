package middleware

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/summerKK/mall-api/pkg/util"
)

//  收集运行时的错误日志并打印
func CollectError(writer io.Writer) gin.HandlerFunc {
	if writer == nil {
		writer = os.Stdout
	}
	colo := color.New(color.FgYellow)

	key := util.GetCtxErrorKey()
	return func(c *gin.Context) {
		// 设置收集错误的collection
		errorList := make([]error, 0, 4)
		c.Set(key, &errorList)

		c.Next()

		value, exists := c.Get(key)
		if exists {
			if errors, ok := value.(*[]error); ok && len(*errors) > 0 {

				buf := bytes.Buffer{}
				for _, err := range *errors {
					buf.WriteString(err.Error())
					buf.WriteByte('\n')
				}

				_, _ = colo.Fprintf(writer, "[%s]\n-------------\n%s-------------\n", time.Now().Format(util.TimeLayout), buf.String())
			}
		}
	}
}
