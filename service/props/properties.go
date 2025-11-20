package props

type Config struct {
	Server    ServerConfig    `yaml:"Server"`
	Logging   LoggingConfig   `yaml:"Logging"`
	Cache     CacheConfig     `yaml:"Cache"`
	Providers ProvidersConfig `yaml:"Providers"`
}

type ServerConfig struct {
	Address     string `yaml:"Address"`
	ContextRoot string `yaml:"ContextRoot"`
	Timeout     int    `yaml:"Timeout"` // seconds
}

type LoggingConfig struct {
	Level       string `yaml:"Level"`       // debug, info, warn, error
	Format      string `yaml:"Format"`      // json
	BodyLogging string `yaml:"BodyLogging"` // none | all | errors
}

type CacheConfig struct {
	TTLSeconds int    `yaml:"TTLSeconds"`
	RedisHost  string `yaml:"RedisHost"`
	RedisPort  string `yaml:"RedisPort"`
}

type ProvidersConfig struct {
	YouTube YoutubeConfig
	Twitch  TwitchConfig
}

type YoutubeConfig struct {
	Provider    string `json:"provider"` // "youtube", "twitch", etc.
	ApiKey      string
	SearchQuery string `json:"searchQuery"`
	ChannelID   string `json:"channelId,omitempty"`
	MaxResults  int    `json:"maxResults,omitempty"`
	Features    []string
	UserID      string
}

type TwitchConfig struct{}
