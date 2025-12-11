package configs

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// viper keys
const (
	EnvKey    = "env"
	Mode      = "mode"
	envPrefix = "skopeoui"
)

var viperInitOnce sync.Once
var vip *viper.Viper

func init() {
	viperInitOnce.Do(func() {
		vip = viper.New()
		vip.SetConfigType("yaml")

		vip.SetEnvPrefix(envPrefix) // will be uppercased automatically
		vip.AutomaticEnv()
		vip.AllowEmptyEnv(true)
		vip.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

		vip.SetDefault(EnvKey, "prod")
		vip.SetDefault(Mode, "prod")

		set_core_default()

		if vip.GetString(Mode) == "debug" {
			go func() {
				logrus.Infoln("pprof is starting...")
				err := http.ListenAndServe(":6060", nil)
				if err != nil {
					logrus.WithError(err).Errorln("pprof server error")
				}
			}()
		}
	})
}

// SetConfigFile explicitly defines the path, name and extension of the config file.
// Viper will use this and not check any of the config paths.
func SetConfigFile(in string) {
	if in != "" {
		vip.SetConfigFile(in)
	} else {
		configName := fmt.Sprintf("config-%s", vip.Get(EnvKey))
		logrus.Debugf("config file is %s", configName+".yaml")
		vip.SetConfigName(configName) // name of config file (without extension)
		file, err := exec.LookPath(os.Args[0])
		if err == nil {
			dir, err := filepath.Abs(file)
			if err == nil {
				dir = filepath.Dir(dir)
				vip.AddConfigPath(dir)
			}
		}
		vip.AddConfigPath(fmt.Sprintf("/%s/", envPrefix))
		vip.AddConfigPath(fmt.Sprintf("/etc/%s/", envPrefix))
		vip.AddConfigPath(fmt.Sprintf("$HOME/.%s", envPrefix))
		vip.AddConfigPath(fmt.Sprintf("/opt/%s/", envPrefix))
		vip.AddConfigPath(".")
		vip.AddConfigPath("./bin")
	}
	err := vip.ReadInConfig()
	if err != nil {
		switch err := err.(type) {
		case viper.ConfigFileNotFoundError:
			logrus.Warnln("No config file found. Using environment variables only.")
		case viper.ConfigParseError:
			logrus.Panicf("Cannot read config file: %s", err)
		default:
			logrus.Warnf("Read config file error: %s", err)
		}
	} else {
		logrus.Infof("Loading config from file %s", vip.ConfigFileUsed())
	}
}

// Set  sets the value for the key in the override register.
func Set(key string, value interface{}) {
	vip.Set(key, value)
}

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
func IsSet(key string) bool {
	return vip.IsSet(key)
}

// GetString returns the value associated with the key as a string.
func GetString(key string) string {
	return vip.GetString(key)
}

// GetInt returns the value associated with the key as a int.
func GetInt(key string) int {
	return vip.GetInt(key)
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool {
	return vip.GetBool(key)
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration {
	return vip.GetDuration(key)
}

func AllSettings() map[string]any {
	return vip.AllSettings()
}

func BindPFlag(key string, flag *pflag.Flag) error {
	return vip.BindPFlag(key, flag)
}

func BindPFlags(flags *pflag.FlagSet) error {
	return vip.BindPFlags(flags)
}

func PrintAllSettings() {
	keys := vip.AllKeys()
	for i, k := range keys {
		v := vip.Get(k)
		if i != (len(keys) - 1) {
			fmt.Printf("%s=%v, ", k, v)
		} else {
			fmt.Printf("%s=%v", k, v)
		}
	}
	fmt.Printf("\n")
}
