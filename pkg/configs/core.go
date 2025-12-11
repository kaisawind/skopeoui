package configs

import "time"

type DB struct {
	Type       string
	Address    string
	Username   string
	Password   string
	Database   string
	Expiration time.Duration
}

var DefaultDB = DB{
	Type:       "bbolt",
	Address:    "./data/",
	Database:   "skopeoui.db",
	Expiration: time.Hour,
}

type Http struct {
	Address string
}

var DefaultHttp = Http{
	Address: ":8080",
}

const (
	HttpAddress = "http.address"

	DBType       = "db.type" // json;mongodb;mysql;sqlite3
	DBAddress    = "db.address"
	DBExpiration = "db.expiration"
	DBDatabase   = "db.database"
	DBUsername   = "db.username"
	DBPassword   = "db.password"
)

func set_core_default() {
	vip.SetDefault(DBType, DefaultDB.Type)
	vip.SetDefault(DBAddress, DefaultDB.Address)
	vip.SetDefault(DBExpiration, DefaultDB.Expiration)
	vip.SetDefault(DBDatabase, DefaultDB.Database)
	vip.SetDefault(DBUsername, DefaultDB.Username)
	vip.SetDefault(DBPassword, DefaultDB.Password)

	vip.SetDefault(HttpAddress, DefaultHttp.Address)
}

func GetDB() DB {
	return DB{
		Type:       vip.GetString(DBType),
		Address:    vip.GetString(DBAddress),
		Expiration: vip.GetDuration(DBExpiration),
		Database:   vip.GetString(DBDatabase),
		Username:   vip.GetString(DBUsername),
		Password:   vip.GetString(DBPassword),
	}
}

func GetHttp() Http {
	return Http{
		Address: vip.GetString(HttpAddress),
	}
}
