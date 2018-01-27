package executor

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/willstudy/redis-executor/pkg/common"
	"github.com/willstudy/redis-executor/pkg/executor"
)

var (
	configFile string

	checkIntervalS int32
	timeoutS       string
	retry          int32
)

var serverCmd = &cobra.Command{
	Use:           "server",
	Short:         "Lanuch redis-executor.",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := initExecutorConfig()
		if err != nil {
			return fmt.Errorf("parse redis executor conf failed: %v", err)
		}

		e, err := executor.New(cfg)
		if err != nil {
			return fmt.Errorf("Create exector failed with %v", err)
		}
		return e.Run()
	},
}

func initEnv() map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envs[pair[0]] = pair[1]
	}
	return envs
}

func initExecutorConfig() (*executor.ExecutorConfig, error) {
	var checkIntervalS, retry int
	var timeoutS string
	var err error

	newlog, err := common.NewLocalSyslogLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to get local syslog: %s", err)
	}
	newlog.SetLevel(log.Level(debugLevel))
	logger := newlog.WithFields(log.Fields{
		"app": "redis-executor",
	})
	envs := initEnv()

	input, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Read file failed with %v", err)
	}
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "checkIntervalS") {
			items := strings.Split(line, "=")
			if len(items) != 2 {
				return nil, fmt.Errorf("checkIntervalS config error, should like key=value")
			}
			checkIntervalS, err = strconv.Atoi(items[1])
			if err != nil {
				return nil, fmt.Errorf("checkIntervalS config error, should be a positive number")
			}
		}
		if strings.HasPrefix(line, "timeoutS") {
			items := strings.Split(line, "=")
			if len(items) != 2 {
				return nil, fmt.Errorf("timeoutS config error, should like key=value")
			}
			timeoutS = items[1]
		}
		if strings.HasPrefix(line, "retry") {
			items := strings.Split(line, "=")
			if len(items) != 2 {
				return nil, fmt.Errorf("retry config error, should like key=value")
			}
			retry, err = strconv.Atoi(items[1])
			if err != nil {
				return nil, fmt.Errorf("retry config error, should be a positive number")
			}
		}
	}
	return &executor.ExecutorConfig{
		CheckIntervalS: checkIntervalS,
		TimeoutS:       timeoutS,
		Retry:          retry,

		Envs:   envs,
		Logger: logger,
	}, nil
}

func init() {
	serverCmd.Flags().StringVar(&configFile, "configFile", "/etc/redis/redis-executor.conf",
		"config file path")

	rootCmd.AddCommand(serverCmd)
}
