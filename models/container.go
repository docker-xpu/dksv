package models

var (
	DefaultInfoLocation = "/var/run/my-docker/%s/"              // 容器基础信息存储路径
	ConfigName          = "config.json"                         // 配置文件名
	ContainerLogFile    = "container.log"                       // 日志文件名
	RootUrl             = "/root/images/"                       // 根目录
	MntUrl              = "/root/mnt/%s/"                       // 挂载路径
	WriteLayerUrl       = "/root/writeLayer/%s/"                // 可写层路径
	DefaultNetworkPath  = "/var/run/my-docker/network/network/" // 默认网络路径
	CgroupPath          = "/sys/fs/cgroup/"                     // cgroup 的根路径
	MyDockerBinPath     = "/root/bin/my-docker"                 // my-docker 可执行文件地址
)

// 容器的资源限制
type memLimit struct {
	Limit string `json:"mem_limit"`
	Used  string `json:"mem_used"`
}
type ContainerResourceConfig struct {
	MemoryLimit memLimit `json:"memoryLimit"` // 内存限制
	CpuShare    string   `json:"cpuShare"`    // CPU 份额(相对权重) 默认为 1024,单个 CPU 为 1024，两个为 2048，以此类推
	CpuSet      string   `json:"cpuSet"`      // CPU 亲和性
}

// 容器的网络资源
type ipConfig struct {
	IP   string `json:"IP"`   // IP
	Mask string `json:"Mask"` // 掩码
}
type ContainerNetWorkConfig struct {
	Name    string   `json:"name"` // 网络名
	IpRange ipConfig `json:"IpRange"`
	Driver  string   `json:"Driver"` // 驱动
}

// 容器信息
type ContainerInfo struct {
	Pid         string                  `json:"pid"`
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Command     string                  `json:"command"`     // 运行命令
	CreateTime  string                  `json:"createTime"`  // 创建时间
	Status      string                  `json:"status"`      // 状态
	Volume      string                  `json:"volume"`      // 容器卷
	PortMapping []string                `json:"portmapping"` // 端口映射
	Limits      ContainerResourceConfig `json:"limits"`      // 资源限制
}
