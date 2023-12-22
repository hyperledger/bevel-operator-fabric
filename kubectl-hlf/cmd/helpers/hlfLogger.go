package helpers

import "github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"

type hlfLogger struct {
}

func (f hlfLogger) Fatal(v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Fatalf(format string, v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Fatalln(v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Panic(v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Panicf(format string, v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Panicln(v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Print(v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Printf(format string, v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Println(v ...interface{}) {
	// do nothing
}

func (f hlfLogger) Debug(args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Debugf(format string, args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Debugln(args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Info(args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Infof(format string, args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Infoln(args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Warn(args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Warnf(format string, args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Warnln(args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Error(args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Errorf(format string, args ...interface{}) {
	// do nothing
}

func (f hlfLogger) Errorln(args ...interface{}) {
	// do nothing
}

type HLFLoggerProvider struct{}

func (f HLFLoggerProvider) GetLogger(module string) api.Logger {
	return hlfLogger{}
}
