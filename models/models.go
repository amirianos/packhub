package models

type RemoteRepository struct {
	Address              string `yaml:"address"`
	CacheDirectory       string `yaml:"cacheDir"`
	CacheExpirationCheck string `yaml:"cacheExpirationCheck"`
	CacheValidTime       int    `yaml:"cacheValidTime"`
}
