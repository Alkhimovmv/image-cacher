package config

type Config struct {
	HTTP  HTTPConfig  `yaml:"http"`
	Img   ImgConfig   `yaml:"img"`
	Cache CacheConfig `yaml:"cache"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

type ImgConfig struct {
	Folder  string `yaml:"folder"`
	Timeout int    `yaml:"timeout"`
}

type CacheConfig struct {
	MaxSize uint   `yaml:"max_size"`
	Path    string `yaml:"path"`
}
