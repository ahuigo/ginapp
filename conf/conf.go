package conf

import (
	"bytes"
	"ginapp/config"
	"ginapp/fslib"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func initGcConf() {
	debug.SetGCPercent(1000)
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("max procs:", runtime.GOMAXPROCS(0), runtime.NumCPU())
}
func init() {
	initGcConf()
	initWkDir()
}

type Config struct {
	Http   Http   `mapstructure:"http"`
	PgConf PgConf `mapstructure:"pg_conf"`
	App    App    `mapstructure:"app"`
}
type App struct {
	Mode string `mapstructure:"mode"`
}
type PgConf struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}
type Http struct {
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

var _config *Config

func GetConf() *Config {
	if _config == nil {
		_loadConfig()
		if err := viper.Unmarshal(&_config); err != nil {
			panic(err)
		}
	}
	return _config
}

// LoadConfig 载入配置
func _loadConfig() {
	var err error
	if fslib.IsValidPath("./config/conf.yaml") {
		viper.SetConfigName("conf")
		viper.AddConfigPath("config/")
		err = viper.ReadInConfig()
	} else if config.ConfEmbedBytes != nil {
		viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
		err = viper.ReadConfig(bytes.NewBuffer(config.ConfEmbedBytes))
	}

	if err != nil {
		log.Fatal("fail to load config file:", err)
	}
}

func initWkDir() {
	if !isInTest() {
		println("in normal mode(not test)")
		return
	}
	println("in test mode")
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func isInTest() bool {
	// strings.HasSuffix(os.Args[0], ".test")
	//Or strings.Contains(os.Args[0], "/_test/")
	//Or return flag.Lookup("test.v") != nil
	return testing.Testing()
}
