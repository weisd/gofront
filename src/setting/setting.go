package setting

import (
	"io/ioutil"
	"os"

	"middleware/session"
	"models"
	"modules/log"
	"modules/pongor"

	"github.com/BurntSushi/toml"
)

type CronService struct {
	Enable     bool
	Schedule   string
	RunAtStart bool
}

type Config struct {
	Web       WebService
	AccessLog AccessLogService

	Logs map[string]map[string]log.LogService

	Models map[string]models.GormService
	Redis  map[string]models.RedisService

	Pongo pongor.PongorOption

	Session session.Options
}

type WebService struct {
	Debug     bool
	Listen    string
	Pprof     string
	Fasthttp  bool
	StaticDir string
}

type AccessLogService struct {
	Enable   bool
	FilePath string
}

var Conf Config

func InitConf(confPath string) (err error) {
	contents, err := ReplaceEnvsFile(confPath)
	if err != nil {
		return err
	}

	if _, err = toml.Decode(contents, &Conf); err != nil {
		return err
	}

	return nil
}

func ReplaceEnvsFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return os.ExpandEnv(string(contents)), nil
}
