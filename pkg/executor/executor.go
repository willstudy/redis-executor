package executor

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type ExecutorConfig struct {
	CheckIntervalS int
	TimeoutS       string
	Retry          int

	Envs   map[string]string
	Master string
	Logger *log.Entry
}

type executor struct {
	checkIntervalS int
	timeoutS       string
	retry          int

	envs   map[string]string
	master string
	logger *log.Entry
}

func New(config *ExecutorConfig) (*executor, error) {
	if err := checkExecutorConfig(config); err != nil {
		return nil, err
	}
	return &executor{
		checkIntervalS: config.CheckIntervalS,
		timeoutS:       config.TimeoutS,
		retry:          config.Retry,

		envs:   config.Envs,
		master: config.Master,
		logger: config.Logger,
	}, nil
}

func (e *executor) Run() error {
	logger := e.logger.WithFields(log.Fields{
		"func": "Run",
	})
	logger.Info("=============================================================")
	logger.Info("===================== START NEW PROCESS =====================")
	var err error
	launch := e.getLanuchType()
	switch launch {
	case MASTER_TYPE:
		err = e.lanuchMaster()
	case SLAVE_TYPE:
		err = e.lanuchSlave()
	case SENTINEL_TYPE:
		err = e.lanuchSentinel()
	default:
		err = fmt.Errorf("Can not setup due to wrong launch type.")
	}
	return err
}

func checkExecutorConfig(config *ExecutorConfig) error {
	if len(config.Envs) == 0 {
		return fmt.Errorf("Init Env failed, should not be empty.")
	}
	if config.Logger == nil {
		config.Logger = log.WithFields(log.Fields{
			"app": "redis-executor",
		})
	}
	return nil
}
