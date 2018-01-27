package executor

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/willstudy/redis-executor/pkg/common"
	"github.com/willstudy/redis-executor/pkg/configs"
)

func (e *executor) getMaster() (string, error) {
	logger := e.logger.WithFields(log.Fields{
		"func": "getMaster",
	})

	cmds := []string{
		configs.TimeCMD,
		e.timeoutS,
		configs.RedisCLI,
		"-h",
		e.envs[configs.RedisSentinelHost],
		"-p",
		strconv.Itoa(configs.RedisSentinelConnPort),
		"-a",
		e.envs[configs.RedisPass],
		"sentinel",
		"get-master-addr-by-name",
		"mymaster",
	}
	logger.Infof("CMD: %v.", cmds)

	ctx := context.TODO()
	for i := 0; i < e.retry; i++ {
		output, err := common.ExecCmd(ctx, cmds)
		if err != nil {
			logger.Warnf("Try %d time, error: %v.", i, err)
		} else {
			logger.Infof("Output of get master by sentinel: %s .", output)
			items := strings.Split(output, "\n")
			logger.Infof("Return master IP: %s.", items[0])
			return items[0], nil
		}
		time.Sleep(time.Duration(e.checkIntervalS) * time.Second)
	}
	return "", fmt.Errorf("Has try %d times and not find master.", e.retry)
}

func (e *executor) getPeer() string {
	logger := e.logger.WithFields(log.Fields{
		"func": "getPeer",
	})

	if hostname, err := os.Hostname(); err != nil {
		panic(err)
	} else {
		logger.Infof("hostname: %s, pod0: %s, pod1: %s", hostname, e.envs[configs.EnvPodDNS0], e.envs[configs.EnvPodDNS1])
		if strings.Contains(e.envs[configs.EnvPodDNS0], hostname) {
			return e.envs[configs.EnvPodDNS1]
		} else {
			return e.envs[configs.EnvPodDNS0]
		}
	}
}

func (e *executor) getHostDNS() string {
	if hostname, err := os.Hostname(); err != nil {
		panic(err)
	} else {
		if strings.Contains(e.envs[configs.EnvPodDNS0], hostname) {
			return e.envs[configs.EnvPodDNS0]
		} else {
			return e.envs[configs.EnvPodDNS1]
		}
	}
}

func (e *executor) isMaster(master string) bool {
	logger := e.logger.WithFields(log.Fields{
		"func": "isMaster",
	})

	cmds := []string{
		configs.TimeCMD,
		e.timeoutS,
		configs.RedisCLI,
		"-h",
		master,
		"-a",
		e.envs[configs.RedisPass],
		"INFO",
	}
	logger.Infof("CMD: %v.", cmds)

	var result bool
	ctx := context.TODO()
	for i := 0; i < e.retry; i++ {
		output, err := common.ExecCmd(ctx, cmds)
		if err != nil {
			logger.Warnf("Try %d time, error: %v.", i, err)
			result = false
		} else {
			logger.Infof("Output of get redis info: %s.", output)
			if strings.Contains(output, "role:master") {
				result = true
			} else {
				result = false
			}
			break
		}
		time.Sleep(time.Duration(e.checkIntervalS) * time.Second)
	}
	return result
}

func (e *executor) redisAlive(address string) bool {
	logger := e.logger.WithFields(log.Fields{
		"func": "masterAlive",
	})

	cmds := []string{
		configs.TimeCMD,
		e.timeoutS,
		configs.RedisCLI,
		"-h",
		address,
		"-a",
		e.envs[configs.RedisPass],
		"PING",
	}
	logger.Infof("CMD: %v.", cmds)

	var result bool
	ctx := context.TODO()
	for i := 0; i < e.retry; i++ {
		if output, err := common.ExecCmd(ctx, cmds); err != nil {
			logger.Warnf("Try %d time, error: %v. \n", i, err)
			result = false
		} else {
			logger.Infof("Output of redis PING: %s. \n", output)
			result = true
			break
		}
		time.Sleep(time.Duration(e.checkIntervalS) * time.Second)
	}
	return result
}
