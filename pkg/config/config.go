package config

import (
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
)

type ConfigReader struct {
	serverConn          string
	dbConn      pgx.ConnPoolConfig
}

func CreateConfigReader() *ConfigReader {
	dbConn_ := pgx.ConnConfig{
		Host: viper.GetString("database.host"),
		Port: uint16(viper.GetUint("database.port")),
		Database: viper.GetString("database.name"),
		User: viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}

	dbConnPool_ := pgx.ConnPoolConfig{
		ConnConfig: dbConn_,
		MaxConnections: viper.GetInt("database.connections"),
	}
	
	return &ConfigReader{
		serverConn:        viper.GetString("server.conn"),
		dbConn: dbConnPool_,
	}
}

func (cr *ConfigReader) GetServerConn() string {
	return cr.serverConn
}


func (cr *ConfigReader) GetDBConn() pgx.ConnPoolConfig {
	return cr.dbConn
}