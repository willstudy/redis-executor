package executor

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/willstudy/redis-executor/pkg/common"
	"github.com/willstudy/redis-executor/pkg/configs"
)

type LanuchType string

const (
	MASTER_TYPE   LanuchType = "master"
	SLAVE_TYPE    LanuchType = "slave"
	SENTINEL_TYPE LanuchType = "sentinel"

	UNKNOWN_TYPE LanuchType = "none"
)

func (e *executor) getLanuchType() LanuchType {
	logger := e.logger.WithFields(log.Fields{
		"func": "getLanuchType",
	})
	// sentinel
	if value, ok := e.envs[configs.RedisSentinel]; ok {
		logger.Infof("get env for sentinel: %s=%s .", configs.RedisSentinel, value)
		if value == "true" {
			return SENTINEL_TYPE
		}
		logger.Warn("get error env for sentinel, should be SENTINEL=true .")
		return UNKNOWN_TYPE
	}

	// for master & slave
	for {
		peer := e.getPeer()
		logger.Infof("Peer is: %s .", peer)
		if m := e.redisAlive(peer); !m {
			logger.Debug("Peer is down, I begin to launch as a master.")
			return MASTER_TYPE
		}

		logger.Debug("Peer is still alive")
		if m := e.isMaster(peer); m {
			logger.Debug("Peer is master, I begin to launch as a slave.")
			e.master = peer
			return SLAVE_TYPE
		} else {
			logger.Debug("Peer is alive, but is not master, waiting for peer to be a master.")
		}
		time.Sleep(time.Duration(e.checkIntervalS) * time.Second)
	}
}

func (e *executor) lanuchMaster() error {
	logger := e.logger.WithFields(log.Fields{
		"func": "lanuchMaster",
	})

	os.MkdirAll(configs.RedisDataDir, os.ModePerm)

	input, err := ioutil.ReadFile(configs.MasterConf)
	if err != nil {
		logger.Errorf("Read file failed with %v", err)
		return err
	}
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "%redis-pass%") {
			lines[i] = "requirepass " + e.envs[configs.RedisPass]
			logger.Infof("Update conf: %s .", lines[i])
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(configs.MasterConf, []byte(output), os.ModePerm)
	if err != nil {
		logger.Errorf("Write file failed with %v", err)
		return err
	}

	cmds := []string{
		configs.RedisCMD,
		configs.MasterConf,
		"--protected-mode no",
	}
	ctx := context.TODO()
	if output, err = common.ExecCmd(ctx, cmds); err != nil {
		logger.Warnf("Launch redis master Output: %v, failed with %v .", output, err)
		return err
	}
	return nil
}

func (e *executor) lanuchSlave() error {
	logger := e.logger.WithFields(log.Fields{
		"func": "lanuchSlave",
	})

	os.MkdirAll(configs.RedisDataDir, os.ModePerm)

	input, err := ioutil.ReadFile(configs.SlaveConf)
	if err != nil {
		logger.Errorf("Read file failed with %v", err)
		return err
	}
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "masterauth") {
			lines[i] = "masterauth " + e.envs[configs.RedisPass]
			logger.Infof("Update conf: %s .", lines[i])
			continue
		}
		if strings.Contains(line, "requirepass") {
			lines[i] = "requirepass " + e.envs[configs.RedisPass]
			logger.Infof("Update conf: %s .", lines[i])
			continue
		}
		if strings.Contains(line, "%master-ip%") {
			lines[i] = "slaveof " + e.master + " 6379"
			logger.Infof("Update conf: %s .", lines[i])
			continue
		}
		if strings.Contains(line, "%slave_ip%") {
			lines[i] = "slave-announce-ip " + e.getHostDNS()
			logger.Infof("Update conf: %s .", lines[i])
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(configs.SlaveConf, []byte(output), os.ModePerm)
	if err != nil {
		logger.Errorf("Write file failed with %v", err)
		return err
	}

	cmds := []string{
		configs.RedisCMD,
		configs.SlaveConf,
		"--protected-mode no",
	}
	ctx := context.TODO()
	if output, err = common.ExecCmd(ctx, cmds); err != nil {
		logger.Warnf("Launch redis slave Output: %v, failed with %v.", output, err)
		return err
	}
	return nil
}

func (e *executor) lanuchSentinel() error {
	logger := e.logger.WithFields(log.Fields{
		"func": "lanuchSentinel",
	})

	for {
		logger.Debug("Try pod0 as the master.")
		if e.isMaster(e.envs[configs.EnvPodDNS0]) {
			e.master = e.envs[configs.EnvPodDNS0]
			logger.Debug("Using pod0 as the master")
			break
		}
		logger.Debug("Try pod1 as the master.")
		if e.isMaster(e.envs[configs.EnvPodDNS1]) {
			e.master = e.envs[configs.EnvPodDNS1]
			logger.Debug("Using pod1 as the master")
			break
		}
		logger.Debugf("No available redis node, wait.")
		time.Sleep(time.Duration(e.checkIntervalS) * time.Second)
	}

	input, err := ioutil.ReadFile(configs.SentinelConf)
	if err != nil {
		logger.Errorf("Read file failed with %v", err)
		return err
	}
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "%password%") {
			lines[i] = "sentinel auth-pass mymaster " + e.envs[configs.RedisPass]
			logger.Infof("Update conf: %s .", lines[i])
			continue
		}
		if strings.Contains(line, "%master%") {
			lines[i] = "sentinel monitor mymaster " + e.master + " 6379 2"
			logger.Infof("Update conf: %s .", lines[i])
		}
		if strings.Contains(line, "%down-failover-time%") {
			lines[i] = "sentinel down-after-milliseconds mymaster " + e.envs[configs.EnvSentinelDownTime]
			logger.Infof("Update conf: %s .", lines[i])
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(configs.SentinelConf, []byte(output), os.ModePerm)
	if err != nil {
		logger.Errorf("Write file failed with %v", err)
		return err
	}

	cmds := []string{
		configs.RedisCMD,
		configs.SentinelConf,
		"--sentinel",
		"--protected-mode no",
	}
	ctx := context.TODO()
	if output, err = common.ExecCmd(ctx, cmds); err != nil {
		logger.Warnf("Launch redis sentinel Output: %v, failed with %v.", output, err)
		return err
	}
	return nil
}
