package db

import (
	"database/sql"
	"fmt"
	"github.com/luren5/filter-base/storage"
	"github.com/luren5/filter-base/utils/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Mysql *sql.DB
)

func init() {
	logrus.Info("mysql:start read viper config")
	viper.SetConfigName("mysql")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Errorf("viper read mysql config error:%v", log.ErrorField(err))
	}
	Mysql = storage.NewMysql()
}

func Exec(sql string, params ...interface{}) uint {
	stm, err := Mysql.Prepare(sql)
	if err != nil {
		log.Error(fmt.Sprintf("执行sql：%s, 预处理失败：%v", sql, err))
		return 0
	}
	defer stm.Close()
	res, err := stm.Exec(params...)
	if err != nil {
		log.Error(fmt.Sprintf("执行语句:%s 失败,参数:%v, 错误信息:%v", sql, params, err))
		return 0
	}
	//log.SInfof(fmt.Sprintf("执行语句:%s, 参数:%v 成功", sql, params))
	id, _ := res.LastInsertId()
	return uint(id)
}

func Exec2(sql string, params ...interface{}) (uint, error) {
	stm, err := Mysql.Prepare(sql)
	if err != nil {
		//log.Error(fmt.Sprintf("执行sql：%s, 预处理失败：%v", sql, err))
		return 0, err
	}
	defer stm.Close()
	res, err := stm.Exec(params...)
	if err != nil {
		//log.Error(fmt.Sprintf("执行语句:%s 失败,参数:%v, 错误信息:%v", sql, params, err))
		return 0, err
	}
	//log.SInfof("执行语句:%s, 参数:%v 成功", sql, params)
	id, _ := res.LastInsertId()
	return uint(id), nil
}
