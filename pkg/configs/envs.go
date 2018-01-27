package configs

const (
	RedisSentinelHost = "REDIS_SENTINEL_SERVICE_HOST"
	RedisSentinel     = "SENTINEL"
	RedisPass         = "REDIS_PASS"

	EnvListenAddr     = "LISTEN_ADDR"
	EnvRedisNamespace = "REDIS_NAMESPACE"

	EnvCPULimit     = "CPU_LIMIT"
	EnvMemoryLimit  = "MEMORY_LIMIT"
	EnvStorageLimit = "STORAGE_LIMIT"

	EnvCRDInstanceName    = "CRD_INSTANCE_NAME"
	EnvServiceAccountName = "SERVICE_ACCOUNT_NAME"
	EnvClusterName        = "CLUSTER_NAME"

	EnvRedisOperatorAddresses     = "REDIS_OPERATOR_ADDRESSES"
	EnvRedisDataOperatorAddresses = "REDIS_DATA_OPERATOR_ADDRESSES"
	EnvRedisPeerServerAddr        = "REDIS_PEER_SERVER_ADDR"
	EnvRedisControllerServerAddr  = "REDIS_CONTROLLER_SERVER_ADDR"
	EnvRedisCheckInterval         = "REDIS_CHECK_INTERVAL"

	EnvRedisRootPassword = "REDIS_ROOT_PASSWORD"
	EnvDeploymentName    = "DEPLOYMENT_NAME"

	EnvPodDNS0 = "REDIS_POD_DNS_0"
	EnvPodDNS1 = "REDIS_POD_DNS_1"

	EnvSentinelDownTime = "SENTINEL_DOWN_TIME"
)

const (
	// For Neutron
	EnvDirectAccessAPIServer = "DIRECT_ACCESS_APISERVER"
	EnvKEUserName            = "KE_USER_NAME"
	EnvKEPasswd              = "KE_PASSWD"
	EnvKEProject             = "KE_PROJECT"
	EnvKEAccountAddr         = "KE_ACCOUNT_ADDR"
	EnvKEK8SAddr             = "KE_K8S_ADDR"

	EnvTLBEnable         = "TLB_ENABLE"
	EnvTLBDescription    = "TLB_DESCRIPTION"
	EnvTLBBandwidthLimit = "TLB_BANDWIDTH_LIMIT"
	EnvTLBChargeMode     = "TLB_CHARGE_MODE"
	EnvTLBIPType         = "TLB_IP_TYPE"
	EnvTLBPolicy         = "TLB_POLICY"
)

const (
	EnvBackupBucket = "BACKUP_BUCKET"
	EnvBackupZone   = "BACKUP_ZONE"
	EnvStorageAK    = "STORAGE_AK"
	EnvStorageSK    = "STORAGE_SK"
)
