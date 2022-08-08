package types

import (
	"time"
)

// Config hold the configuration of these software
type Config struct {
	MinVersion            string
	EnableSyslog          bool
	AirflowMesosScheduler string
	LogLevel              string
	AppName               string
	Wait                  string
	AWSWait               string
	SSL                   bool
	SkipSSL               bool
	PollInterval          time.Duration
	PollTimeout           time.Duration
	WaitTimeout           time.Duration
	AWSTerminateWait      time.Duration
	RedisServer           string
	RedisPassword         string
	RedisDB               int
	RedisPrefix           string
	APIUsername           string
	APIPassword           string
	AWSRegion             string
	AWSInstance16         string
	AWSInstance32         string
	AWSInstance64         string
	AWSLaunchTemplateID   string
	MesosAgentUsername    string
	MesosAgentPassword    string
	MesosAgentPort        string
	MesosAgentSSL         bool
}

type DagTask struct {
	DagID         string `json:"dag_id"`
	TaskID        string `json:"task_id"`
	RunID         string `json:"run_id"`
	TryNumber     int    `json:"try_number"`
	ASG           bool   `json:"ASG" default:"false"`
	MesosExecutor struct {
		Cpus     float64 `json:"cpus"`
		MemLimit int     `json:"mem_limit" default:"2048"`
	} `json:"MesosExecutor"`
}
type MesosAgentState struct {
	Version      string   `json:"version"`
	GitSha       string   `json:"git_sha"`
	GitTag       string   `json:"git_tag"`
	BuildDate    string   `json:"build_date"`
	BuildTime    int      `json:"build_time"`
	BuildUser    string   `json:"build_user"`
	StartTime    float64  `json:"start_time"`
	ID           string   `json:"id"`
	Pid          string   `json:"pid"`
	Hostname     string   `json:"hostname"`
	Capabilities []string `json:"capabilities"`
	Resources    struct {
		Disk  int    `json:"disk"`
		Mem   int    `json:"mem"`
		Gpus  int    `json:"gpus"`
		Cpus  int    `json:"cpus"`
		Ports string `json:"ports"`
	} `json:"resources"`
	ReservedResources struct {
	} `json:"reserved_resources"`
	UnreservedResources struct {
		Disk  int    `json:"disk"`
		Mem   int    `json:"mem"`
		Gpus  int    `json:"gpus"`
		Cpus  int    `json:"cpus"`
		Ports string `json:"ports"`
	} `json:"unreserved_resources"`
	ReservedResourcesFull struct {
	} `json:"reserved_resources_full"`
	UnreservedResourcesFull []struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Scalar struct {
			Value int `json:"value"`
		} `json:"scalar,omitempty"`
		Role   string `json:"role"`
		Ranges struct {
			Range []struct {
				Begin int `json:"begin"`
				End   int `json:"end"`
			} `json:"range"`
		} `json:"ranges,omitempty"`
	} `json:"unreserved_resources_full"`
	ReservedResourcesAllocated struct {
	} `json:"reserved_resources_allocated"`
	UnreservedResourcesAllocated struct {
		Disk  int     `json:"disk"`
		Mem   int     `json:"mem"`
		Gpus  int     `json:"gpus"`
		Cpus  float64 `json:"cpus"`
		Ports string  `json:"ports"`
	} `json:"unreserved_resources_allocated"`
	ResourceProviders []interface{} `json:"resource_providers"`
	Attributes        struct {
	} `json:"attributes"`
	MasterHostname string `json:"master_hostname"`
	LogDir         string `json:"log_dir"`
	Flags          struct {
		AppcSimpleDiscoveryURIPrefix      string `json:"appc_simple_discovery_uri_prefix"`
		AppcStoreDir                      string `json:"appc_store_dir"`
		AuthenticateHTTPReadonly          string `json:"authenticate_http_readonly"`
		AuthenticateHTTPReadwrite         string `json:"authenticate_http_readwrite"`
		Authenticatee                     string `json:"authenticatee"`
		AuthenticationBackoffFactor       string `json:"authentication_backoff_factor"`
		AuthenticationTimeoutMax          string `json:"authentication_timeout_max"`
		AuthenticationTimeoutMin          string `json:"authentication_timeout_min"`
		Authorizer                        string `json:"authorizer"`
		CgroupsCPUEnablePidsAndTidsCount  string `json:"cgroups_cpu_enable_pids_and_tids_count"`
		CgroupsDestroyTimeout             string `json:"cgroups_destroy_timeout"`
		CgroupsEnableCfs                  string `json:"cgroups_enable_cfs"`
		CgroupsHierarchy                  string `json:"cgroups_hierarchy"`
		CgroupsLimitSwap                  string `json:"cgroups_limit_swap"`
		CgroupsRoot                       string `json:"cgroups_root"`
		ContainerDiskWatchInterval        string `json:"container_disk_watch_interval"`
		Containerizers                    string `json:"containerizers"`
		DefaultRole                       string `json:"default_role"`
		DisallowSharingAgentIpcNamespace  string `json:"disallow_sharing_agent_ipc_namespace"`
		DisallowSharingAgentPidNamespace  string `json:"disallow_sharing_agent_pid_namespace"`
		DiskProfileAdaptor                string `json:"disk_profile_adaptor"`
		DiskWatchInterval                 string `json:"disk_watch_interval"`
		Docker                            string `json:"docker"`
		DockerIgnoreRuntime               string `json:"docker_ignore_runtime"`
		DockerKillOrphans                 string `json:"docker_kill_orphans"`
		DockerRegistry                    string `json:"docker_registry"`
		DockerRemoveDelay                 string `json:"docker_remove_delay"`
		DockerSocket                      string `json:"docker_socket"`
		DockerStopTimeout                 string `json:"docker_stop_timeout"`
		DockerStoreDir                    string `json:"docker_store_dir"`
		DockerVolumeCheckpointDir         string `json:"docker_volume_checkpoint_dir"`
		DockerVolumeChown                 string `json:"docker_volume_chown"`
		EnforceContainerDiskQuota         string `json:"enforce_container_disk_quota"`
		ExecutorRegistrationTimeout       string `json:"executor_registration_timeout"`
		ExecutorReregistrationTimeout     string `json:"executor_reregistration_timeout"`
		ExecutorShutdownGracePeriod       string `json:"executor_shutdown_grace_period"`
		FetcherCacheDir                   string `json:"fetcher_cache_dir"`
		FetcherCacheSize                  string `json:"fetcher_cache_size"`
		FetcherStallTimeout               string `json:"fetcher_stall_timeout"`
		FrameworksHome                    string `json:"frameworks_home"`
		GcDelay                           string `json:"gc_delay"`
		GcDiskHeadroom                    string `json:"gc_disk_headroom"`
		GcNonExecutorContainerSandboxes   string `json:"gc_non_executor_container_sandboxes"`
		Help                              string `json:"help"`
		Hostname                          string `json:"hostname"`
		HostnameLookup                    string `json:"hostname_lookup"`
		HTTPCommandExecutor               string `json:"http_command_executor"`
		HTTPExecutorDomainSockets         string `json:"http_executor_domain_sockets"`
		HTTPHeartbeatInterval             string `json:"http_heartbeat_interval"`
		ImageProviders                    string `json:"image_providers"`
		ImageProvisionerBackend           string `json:"image_provisioner_backend"`
		InitializeDriverLogging           string `json:"initialize_driver_logging"`
		Isolation                         string `json:"isolation"`
		Launcher                          string `json:"launcher"`
		LauncherDir                       string `json:"launcher_dir"`
		LogDir                            string `json:"log_dir"`
		Logbufsecs                        string `json:"logbufsecs"`
		LoggingLevel                      string `json:"logging_level"`
		Master                            string `json:"master"`
		MaxCompletedExecutorsPerFramework string `json:"max_completed_executors_per_framework"`
		MemoryProfiling                   string `json:"memory_profiling"`
		ModulesDir                        string `json:"modules_dir"`
		NetworkCniConfigDir               string `json:"network_cni_config_dir"`
		NetworkCniMetrics                 string `json:"network_cni_metrics"`
		NetworkCniPluginsDir              string `json:"network_cni_plugins_dir"`
		NetworkCniRootDirPersist          string `json:"network_cni_root_dir_persist"`
		OversubscribedResourcesInterval   string `json:"oversubscribed_resources_interval"`
		PerfDuration                      string `json:"perf_duration"`
		PerfInterval                      string `json:"perf_interval"`
		Port                              string `json:"port"`
		QosCorrectionIntervalMin          string `json:"qos_correction_interval_min"`
		Quiet                             string `json:"quiet"`
		ReconfigurationPolicy             string `json:"reconfiguration_policy"`
		Recover                           string `json:"recover"`
		RecoveryTimeout                   string `json:"recovery_timeout"`
		RegistrationBackoffFactor         string `json:"registration_backoff_factor"`
		ResourceProviderConfigDir         string `json:"resource_provider_config_dir"`
		RevocableCPULowPriority           string `json:"revocable_cpu_low_priority"`
		RuntimeDir                        string `json:"runtime_dir"`
		SandboxDirectory                  string `json:"sandbox_directory"`
		Strict                            string `json:"strict"`
		SwitchUser                        string `json:"switch_user"`
		SystemdEnableSupport              string `json:"systemd_enable_support"`
		SystemdRuntimeDirectory           string `json:"systemd_runtime_directory"`
		Version                           string `json:"version"`
		WorkDir                           string `json:"work_dir"`
		ZkSessionTimeout                  string `json:"zk_session_timeout"`
	} `json:"flags"`
	Frameworks []struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		User            string `json:"user"`
		FailoverTimeout int    `json:"failover_timeout"`
		Checkpoint      bool   `json:"checkpoint"`
		Hostname        string `json:"hostname"`
		Principal       string `json:"principal"`
		Role            string `json:"role"`
		Executors       []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Source    string `json:"source"`
			Container string `json:"container"`
			Directory string `json:"directory"`
			Resources struct {
				Disk  int     `json:"disk"`
				Mem   int     `json:"mem"`
				Gpus  int     `json:"gpus"`
				Cpus  float64 `json:"cpus"`
				Ports string  `json:"ports"`
			} `json:"resources"`
			Role  string `json:"role"`
			Tasks []struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				FrameworkID string `json:"framework_id"`
				ExecutorID  string `json:"executor_id"`
				SlaveID     string `json:"slave_id"`
				State       string `json:"state"`
				Resources   struct {
					Disk  int     `json:"disk"`
					Mem   int     `json:"mem"`
					Gpus  int     `json:"gpus"`
					Cpus  float64 `json:"cpus"`
					Ports string  `json:"ports"`
				} `json:"resources"`
				Role     string `json:"role"`
				Statuses []struct {
					State           string  `json:"state"`
					Timestamp       float64 `json:"timestamp"`
					ContainerStatus struct {
						ContainerID struct {
							Value string `json:"value"`
						} `json:"container_id"`
						NetworkInfos []struct {
							IPAddresses []struct {
								Protocol  string `json:"protocol"`
								IPAddress string `json:"ip_address"`
							} `json:"ip_addresses"`
						} `json:"network_infos"`
					} `json:"container_status"`
				} `json:"statuses"`
				Discovery struct {
					Visibility string `json:"visibility"`
					Name       string `json:"name"`
					Ports      struct {
						Ports []struct {
							Number   int    `json:"number"`
							Name     string `json:"name"`
							Protocol string `json:"protocol"`
						} `json:"ports"`
					} `json:"ports"`
				} `json:"discovery"`
				Container struct {
					Type     string `json:"type"`
					Hostname string `json:"hostname"`
					Docker   struct {
						Image        string `json:"image"`
						Network      string `json:"network"`
						PortMappings []struct {
							HostPort      int    `json:"host_port"`
							ContainerPort int    `json:"container_port"`
							Protocol      string `json:"protocol"`
						} `json:"port_mappings"`
						Privileged bool `json:"privileged"`
						Parameters []struct {
							Key   string `json:"key"`
							Value string `json:"value"`
						} `json:"parameters"`
						ForcePullImage bool `json:"force_pull_image"`
					} `json:"docker"`
					NetworkInfos []struct {
						Name string `json:"name"`
					} `json:"network_infos"`
				} `json:"container"`
			} `json:"tasks"`
			QueuedTasks    []interface{} `json:"queued_tasks"`
			CompletedTasks []interface{} `json:"completed_tasks"`
		} `json:"executors"`
		CompletedExecutors []interface{} `json:"completed_executors"`
	} `json:"frameworks"`
	CompletedFrameworks []struct {
		ID                 string        `json:"id"`
		Name               string        `json:"name"`
		User               string        `json:"user"`
		FailoverTimeout    int           `json:"failover_timeout"`
		Checkpoint         bool          `json:"checkpoint"`
		Hostname           string        `json:"hostname"`
		Principal          string        `json:"principal,omitempty"`
		Role               string        `json:"role"`
		Executors          []interface{} `json:"executors"`
		CompletedExecutors []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Source    string `json:"source"`
			Container string `json:"container"`
			Directory string `json:"directory"`
			Resources struct {
				Disk int     `json:"disk"`
				Mem  int     `json:"mem"`
				Gpus int     `json:"gpus"`
				Cpus float64 `json:"cpus"`
			} `json:"resources"`
			Role   string `json:"role"`
			Labels []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"labels"`
			Tasks          []interface{} `json:"tasks"`
			QueuedTasks    []interface{} `json:"queued_tasks"`
			CompletedTasks []struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				FrameworkID string `json:"framework_id"`
				ExecutorID  string `json:"executor_id"`
				SlaveID     string `json:"slave_id"`
				State       string `json:"state"`
				Resources   struct {
					Disk int     `json:"disk"`
					Mem  int     `json:"mem"`
					Gpus int     `json:"gpus"`
					Cpus float64 `json:"cpus"`
				} `json:"resources"`
				Role     string `json:"role"`
				Statuses []struct {
					State           string  `json:"state"`
					Timestamp       float64 `json:"timestamp"`
					ContainerStatus struct {
						ContainerID struct {
							Value string `json:"value"`
						} `json:"container_id"`
						NetworkInfos []struct {
							IPAddresses []struct {
								Protocol  string `json:"protocol"`
								IPAddress string `json:"ip_address"`
							} `json:"ip_addresses"`
						} `json:"network_infos"`
					} `json:"container_status"`
				} `json:"statuses"`
				Labels []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"labels"`
				Discovery struct {
					Visibility string `json:"visibility"`
					Name       string `json:"name"`
					Ports      struct {
					} `json:"ports"`
				} `json:"discovery"`
				Container struct {
					Type     string `json:"type"`
					Hostname string `json:"hostname"`
					Docker   struct {
						Image      string `json:"image"`
						Network    string `json:"network"`
						Privileged bool   `json:"privileged"`
						Parameters []struct {
							Key   string `json:"key"`
							Value string `json:"value"`
						} `json:"parameters"`
						ForcePullImage bool `json:"force_pull_image"`
					} `json:"docker"`
					NetworkInfos []struct {
						Name string `json:"name"`
					} `json:"network_infos"`
					LinuxInfo struct {
						ShareCgroups bool `json:"share_cgroups"`
					} `json:"linux_info"`
				} `json:"container"`
			} `json:"completed_tasks"`
		} `json:"completed_executors"`
	} `json:"completed_frameworks"`
}
