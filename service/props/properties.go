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
	Provider string `json:"provider"` // "youtube"
	Enabled  bool
	client   Api
	ApiKey   string
	Features YouTubeFeatures `json:"features"`
}
type YouTubeFeatures struct {
	Trending         *TrendingConfig         `json:"trending,omitempty"`
	ContinueWatching *ContinueWatchingConfig `json:"continueWatching,omitempty"`
	// Add more in the future...
}

type TrendingConfig struct {
	MaxResults int    `json:"maxResults,omitempty"`
	Region     string `json:"region,omitempty"`
}

type ContinueWatchingConfig struct {
	MaxResults int `json:"maxResults,omitempty"`
}

type TwitchConfig struct {
	Provider     string `json:"provider"` // "twitch"
	Enabled      bool
	Client ClientProps  `json:"client"`
	ClientID     string `json:"clientId"`     // Twitch Client-ID
	ClientSecret string `json:"clientSecret"` // Twitch Client-Secret
}

// Api holds api configuration for any service
type Api struct {
	Name       string               `yaml:"name" envconfig:"NAME"`
	Host       string               `yaml:"host" envconfig:"DNS_NAME"`
	BasePath   string               `yaml:"base-path" envconfig:"BASE_PATH"`
	Scheme     string               `yaml:"scheme" envconfig:"SCHEME"`
	Operations map[string]Operation `yaml:"operations"`
	Client     ClientProps          `yaml:"client" envconfig:"CLIENT"`
}

// Operation represents various end-point/resources of a client API
type Operation struct {
	Method      string `yaml:"method" envconfig:"METHOD"`
	PathPattern string `yaml:"path" envconfig:"PATH_PATTERN"`
}
