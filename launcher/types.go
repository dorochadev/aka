package launcher

type LauncherType string

const (
	TypeApplication LauncherType = "app"
	TypeURL         LauncherType = "url"
	TypeSSH         LauncherType = "ssh"
	TypeCommand     LauncherType = "cmd"
)

type LauncherMetadata struct {
	Type      LauncherType      `json:"type"`
	Target    string            `json:"target"`
	Env       map[string]string `json:"env,omitempty"`
	SSHConfig *SSHConfig        `json:"ssh_config,omitempty"`
}

type SSHConfig struct {
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
}
