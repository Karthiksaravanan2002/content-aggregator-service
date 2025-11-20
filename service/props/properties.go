package props

type Config struct {
	Server    ServerConfig    `yaml:"server" json:"Server"`
	Logging   LoggingConfig   `yaml:"logging" json:"Logging"`
	Cache     CacheConfig     `yaml:"cache" json:"Cache"`
	Providers ProvidersConfig `yaml:"providers" json:"Providers"`
}

type ServerConfig struct {
	Address     string `yaml:"address" json:"Address"`
	ContextRoot string `yaml:"contextRoot" json:"ContextRoot"`
	Timeout     int    `yaml:"timeout" json:"Timeout"` 
}

type LoggingConfig struct {
	Level       string `yaml:"level" json:"Level"`             // debug, info, warn, error
	Format      string `yaml:"format" json:"Format"`           // json
	BodyLogging string `yaml:"bodyLogging" json:"BodyLogging"` // none | all | errors
}

type CacheConfig struct {
	TTLSeconds int    `yaml:"ttlSeconds" json:"TTLSeconds"`
	RedisHost  string `yaml:"redisHost" json:"RedisHost"`
	RedisPort  string `yaml:"redisPort" json:"RedisPort"`
}

type ProvidersConfig struct {
	YouTube YouTubeConfig `yaml:"youTube" json:"YouTube"`
	Twitch  TwitchConfig  `yaml:"twitch" json:"Twitch"`
}

type YouTubeConfig struct {
	Provider    string `json:"provider"` // "youtube", "twitch", etc.
	ApiKey      string
	SearchQuery string `json:"searchQuery"`
	ChannelID   string `json:"channelId,omitempty"`
	MaxResults  int    `json:"maxResults,omitempty"`
	Features    []string
	UserID      string
}

type TwitchConfig struct {
	ClientID     string `yaml:"clientID" json:"ClientID"`
	ClientSecret string `yaml:"clientSecret" json:"ClientSecret"`
	RedirectURI  string `yaml:"redirectURI" json:"RedirectURI"`
}
