package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

// trustedProxies are the only peers whose X-Forwarded-For the server believes:
// nginx on the same machine. Nobody else connects from loopback, so a client
// can't fake its IP; a direct connection keeps its own address.
var trustedProxies = []string{"127.0.0.1", "::1"}

func NewServer(origins ...string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(logFormatter), gin.Recovery())

	// nginx sets X-Forwarded-For to the client's address (not appended to), so
	// only that header is read.
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		panic(err)
	}
	router.RemoteIPHeaders = []string{"X-Forwarded-For"}

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = append(corsCfg.AllowOrigins, "https://"+models.DefaultDomain)
	corsCfg.AllowOrigins = append(corsCfg.AllowOrigins, "https://"+models.DevDomain)
	corsCfg.AllowOrigins = append(corsCfg.AllowOrigins, origins...)
	corsCfg.AllowCredentials = true

	router.Use(cors.New(corsCfg))

	return router
}

// logFormatter is gin's default request line plus the signed-in user's ID
// ("-" for a visitor who isn't signed in), so a user's requests can be told
// apart in the log. Only the ID: a name would cost a database query per
// request, and the ID is enough to look the user up.
func logFormatter(param gin.LogFormatterParams) string {
	user := "-"
	if id, ok := param.Keys[reqctx.UserIDKey].(int); ok && id != reqctx.ZeroUserID {
		user = strconv.Itoa(id)
	}

	return fmt.Sprintf("[GIN] %v | %3d | %8v | %15s | user %-5s | %-7s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.StatusCode,
		param.Latency.Round(time.Microsecond),
		param.ClientIP,
		user,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}
