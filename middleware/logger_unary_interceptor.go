package middleware

import (
	configs "github.com/Misfit-SW-China/cloud-fitness/configs"
	"github.com/Sirupsen/logrus"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

var accessLog *logrus.Logger

func init() {
	filename, _ := filepath.Abs("./logs/" + configs.AppConfig.GOENV + ".log")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	accessLog = logrus.New()
	accessLog.Formatter = new(logrus.JSONFormatter)
	accessLog.Out = file
	accessLog.Level = logrus.InfoLevel
	logrus.NewEntry(accessLog)
}

func LoggerStreamInterceptor(
	srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) (err error) {
	defer log(ss.Context(), info, time.Now().UTC(), err)
	return handler(srv, ss)
}

func log(ctx context.Context, info *grpc.StreamServerInfo, start time.Time, grpcErr error) {
	code := grpc.Code(grpcErr)
	desc := grpc.ErrorDesc(grpcErr)

	service := "unknowService"
	method := "unknowMethod"

	fullMethodSlice := strings.Split(info.FullMethod, "/")
	if len(fullMethodSlice) == 3 {
		service = fullMethodSlice[1]
		method = fullMethodSlice[2]
	}

	duration := time.Since(start)
	var addr string
	var tls uint16
	var cipher uint16
	var network string
	if pr, ok := peer.FromContext(ctx); ok {
		if pr.AuthInfo != nil {
			if info, ok := pr.AuthInfo.(credentials.TLSInfo); ok {
				tls = info.State.Version
				cipher = info.State.CipherSuite
			}
		}

		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}

		network = pr.Addr.Network()
	}

	var level string
	switch code {
	case codes.OK:
		level = "info"
	case codes.Unknown, codes.Internal:
		level = "error"
	default:
		level = "warn"
	}

	accessLog.WithFields(logrus.Fields{
		"time":     time.Now().Format(time.RFC3339),
		"service":  service,
		"method":   method,
		"status":   code.String(),
		"code":     uint32(code),
		"duration": duration.Nanoseconds() / int64(time.Millisecond),
		"error":    desc,
		"level":    level,
		"remote":   addr,
		"tls":      tls,
		"cipher":   cipher,
		"network":  network,
	}).Info("Fitness")
}
