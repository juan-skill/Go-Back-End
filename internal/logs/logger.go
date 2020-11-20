package logs

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

// InitLogger function will initialize a specific development logger
func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	sugar = logger.Sugar()

	return nil
}

// Log function will return a instance of logger
func Log() *zap.SugaredLogger {
	return sugar
}
