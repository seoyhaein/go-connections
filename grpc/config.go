package grpc

// TODO toml 안쓸계획임. 그냥 yaml 사용할 예정임.
// Config provides containerd configuration data for the server
type Config struct {
	Debug bool `toml:"debug"`

	// Root is the path to a directory where buildkit will store persistent data
	Root string `toml:"root"`

	// Entitlements e.g. security.insecure, network.host
	Entitlements []string `toml:"insecure-entitlements"`
	// GRPC configuration settings
	GRPC GRPCConfig `toml:"grpc"`
}

type GRPCConfig struct {
	Address      []string `toml:"address"`
	DebugAddress string   `toml:"debugAddress"`
	UID          *int     `toml:"uid"`
	GID          *int     `toml:"gid"`

	TLS TLSConfig `toml:"tls"`
	// MaxRecvMsgSize int    `toml:"max_recv_message_size"`
	// MaxSendMsgSize int    `toml:"max_send_message_size"`
}

type TLSConfig struct {
	Cert string `toml:"cert"`
	Key  string `toml:"key"`
	CA   string `toml:"ca"`
}
