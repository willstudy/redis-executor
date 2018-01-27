package configs

const (
	RedisDataDir    = "/redis-data"
	RedisDataVolume = "redis-data-volume"

	MasterConf   = "/etc/redis/redis-master.conf"
	SlaveConf    = "/etc/redis/redis-slave.conf"
	SentinelConf = "/etc/redis/redis-sentinel.conf"

	RedisCLI      = "/usr/local/bin/redis-cli"
	SupervisorCLI = "/usr/bin/supervisorctl"
	TimeCMD       = "/usr/bin/timeout"
	RedisCMD      = "/usr/local/bin/redis-server"

	HomeDir    = "/home/redis"
	BackupDir  = "/home/redis/backup"
	BackupCMD  = "/usr/bin/rdiff-backup"
	RestoreCMD = "/usr/bin/rdiff-backup"
)
