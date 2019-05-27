package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var wdb *sql.DB

type ConnectConfig struct {
	Username string
	Password string
	Host     string
	Port     uint
	DB       string
	Charset  string
}

//var cc = &ConnectConfig{
//	Username: "root",
//	Password: "123456",
//	Host:     "192.168.1.164",
//	Port:     3306,
//	DB:       "filter",
//	Charset:  "utf8",
//}

func getCC() (*ConnectConfig) {
	var cc *ConnectConfig
	if viper.IsSet("mysql") {
		mysql := viper.GetStringMap("mysql")
		cc = &ConnectConfig{}
		cc.Username = mysql["user"].(string)
		cc.Password = mysql["pass"].(string)
		cc.Host = mysql["host"].(string)
		cc.Port = uint(mysql["port"].(int))
		cc.DB = mysql["db"].(string)
		cc.Charset = "utf8"
		fmt.Println("use viper configed mysql")
		logrus.Info("use viper configed mysql")
	} else {
		cc = &ConnectConfig{
			Username: "root",
			Password: "123456",
			Host:     "192.168.1.164",
			Port:     3306,
			DB:       "filter",
			Charset:  "utf8",
		}
		logrus.Info("use default configed mysql")
	}
	return cc
}

func NewMysql() *sql.DB {
	if wdb == nil {
		if d, err := newMysqlDB(getCC()); err != nil {
			panic(err)
		} else {
			wdb = d
		}
	}
	return wdb
}

func NewMysqlWithConf(conf *ConnectConfig) *sql.DB {
	d, err := newMysqlDB(conf)
	if err != nil {
		panic(err)
	}
	return d
}

func newMysqlDB(c *ConnectConfig) (*sql.DB, error) {
	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", c.Username, c.Password, c.Host, c.Port, c.DB, c.Charset)
	return sql.Open("mysql", connInfo)
}
