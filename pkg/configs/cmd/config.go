package cmd

import "github.com/kaisawind/skopeoui/pkg/configs"

var dbEnvs = map[string]*env{
	configs.DBType: {
		Usage:   "database type",
		Default: configs.DefaultDB.Type,
	},
	configs.DBAddress: {
		Usage:   "database address",
		Default: configs.DefaultDB.Address,
	},
	configs.DBDatabase: {
		Usage:   "database name",
		Default: configs.DefaultDB.Database,
	},
	configs.DBUsername: {
		Usage:   "database username",
		Default: configs.DefaultDB.Username,
	},
	configs.DBPassword: {
		Usage:   "database password",
		Default: configs.DefaultDB.Password,
	},
	configs.DBExpiration: {
		Usage:   "database expiration",
		Default: configs.DefaultDB.Expiration.String(),
	},
}

var httpEnvs = map[string]*env{
	configs.HttpAddress: {
		Usage:   "http address",
		Default: configs.DefaultHttp.Address,
	},
}

var envs = map[string]*env{}

func DBEnvs() map[string]*env {
	return dbEnvs
}

func HttpEnvs() map[string]*env {
	return httpEnvs
}

func Envs() map[string]*env {
	for k, v := range dbEnvs {
		envs[k] = v
	}
	for k, v := range httpEnvs {
		envs[k] = v
	}
	return envs
}
