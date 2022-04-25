package mesos

type MesosTasks struct {
	Tasks []MesosTask `json:"tasks"`
}

type MesosTask struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FrameworkID string `json:"framework_id"`
	ExecutorID  string `json:"executor_id"`
	SlaveID     string `json:"slave_id"`
	State       string `json:"state"`
	Resources   struct {
		Disk  float64 `json:"disk"`
		Mem   float64 `json:"mem"`
		Gpus  float64 `json:"gpus"`
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
		Healthy bool `json:"healthy,omitempty"`
	} `json:"statuses"`
	Labels    []MesosLabels `json:"labels"`
	Discovery struct {
		Visibility string `json:"visibility"`
		Name       string `json:"name"`
		Ports      struct {
			Ports []MesosPorts `json:"ports"`
		} `json:"ports"`
	} `json:"discovery"`
	Container struct {
		Type   string `json:"type"`
		Docker struct {
			Image        string `json:"image"`
			Network      string `json:"network"`
			PortMappings []struct {
				HostPort      uint32 `json:"host_port"`
				ContainerPort uint32 `json:"container_port"`
				Protocol      string `json:"protocol"`
			} `json:"port_mappings"`
			Privileged bool `json:"privileged"`
			Parameters []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"parameters"`
			ForcePullImage bool `json:"force_pull_image"`
		} `json:"docker"`
	} `json:"container"`
	HealthCheck struct {
		DelaySeconds        float64 `json:"delay_seconds"`
		IntervalSeconds     float64 `json:"interval_seconds"`
		TimeoutSeconds      float64 `json:"timeout_seconds"`
		ConsecutiveFailures float64 `json:"consecutive_failures"`
		GracePeriodSeconds  float64 `json:"grace_period_seconds"`
		Type                string  `json:"type"`
		HTTP                struct {
			Protocol string `json:"protocol"`
			Scheme   string `json:"scheme"`
			Port     int    `json:"port"`
			Path     string `json:"path"`
		} `json:"http"`
	} `json:"health_check"`
}

type MesosLabels struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MesosPorts struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Labels   struct {
		Labels []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"labels"`
	} `json:"labels"`
}

type MesosAgent struct {
	Slaves []struct {
		ID         string `json:"id"`
		Hostname   string `json:"hostname"`
		Port       int    `json:"port"`
		Attributes struct {
		} `json:"attributes"`
		Pid              string  `json:"pid"`
		RegisteredTime   float64 `json:"registered_time"`
		ReregisteredTime float64 `json:"reregistered_time"`
		Resources        struct {
			Disk  float64 `json:"disk"`
			Mem   float64 `json:"mem"`
			Gpus  float64 `json:"gpus"`
			Cpus  float64 `json:"cpus"`
			Ports string  `json:"ports"`
		} `json:"resources"`
		UsedResources struct {
			Disk  float64 `json:"disk"`
			Mem   float64 `json:"mem"`
			Gpus  float64 `json:"gpus"`
			Cpus  float64 `json:"cpus"`
			Ports string  `json:"ports"`
		} `json:"used_resources"`
		OfferedResources struct {
			Disk float64 `json:"disk"`
			Mem  float64 `json:"mem"`
			Gpus float64 `json:"gpus"`
			Cpus float64 `json:"cpus"`
		} `json:"offered_resources"`
		ReservedResources struct {
		} `json:"reserved_resources"`
		UnreservedResources struct {
			Disk  float64 `json:"disk"`
			Mem   float64 `json:"mem"`
			Gpus  float64 `json:"gpus"`
			Cpus  float64 `json:"cpus"`
			Ports string  `json:"ports"`
		} `json:"unreserved_resources"`
		Active                bool     `json:"active"`
		Deactivated           bool     `json:"deactivated"`
		Version               string   `json:"version"`
		Capabilities          []string `json:"capabilities"`
		ReservedResourcesFull struct {
		} `json:"reserved_resources_full"`
		UnreservedResourcesFull []struct {
			Name   string `json:"name"`
			Type   string `json:"type"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar,omitempty"`
			Role   string `json:"role"`
			Ranges struct {
				Range []struct {
					Begin int `json:"begin"`
					End   int `json:"end"`
				} `json:"range"`
			} `json:"ranges,omitempty"`
		} `json:"unreserved_resources_full"`
		UsedResourcesFull []struct {
			Name   string `json:"name"`
			Type   string `json:"type"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar,omitempty"`
			Role           string `json:"role"`
			AllocationInfo struct {
				Role string `json:"role"`
			} `json:"allocation_info"`
			Ranges struct {
				Range []struct {
					Begin int `json:"begin"`
					End   int `json:"end"`
				} `json:"range"`
			} `json:"ranges,omitempty"`
		} `json:"used_resources_full"`
		OfferedResourcesFull []interface{} `json:"offered_resources_full"`
	} `json:"slaves"`
	RecoveredSlaves []interface{} `json:"recovered_slaves"`
}

type MesosAgentContainers []struct {
	ContainerID  string `json:"container_id"`
	ExecutorID   string `json:"executor_id"`
	ExecutorName string `json:"executor_name"`
	FrameworkID  string `json:"framework_id"`
	Source       string `json:"source"`
	Status       struct {
		ContainerID struct {
			Value string `json:"value"`
		} `json:"container_id"`
	} `json:"status"`
}
