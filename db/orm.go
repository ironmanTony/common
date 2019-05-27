package db

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

const TAG_ORM_NAME = "db"

func Insert(ptr interface{}, tableName string) error {
	v := reflect.ValueOf(ptr).Elem()
	var values []interface{}
	var fieldsName []string
	var placeHolders []string
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		label := fieldInfo.Tag.Get(TAG_ORM_NAME)
		if label != "" {
			fieldsName = append(fieldsName, label)
			values = append(values, getValue(v.Field(i)))
			placeHolders = append(placeHolders, "?")
		}
	}
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", tableName, strings.Join(fieldsName, ","), strings.Join(placeHolders, ","))
	_, err := Exec2(sql, values...)
	return err
}

func InsertSlice(tableName string, ignoreDuplicate bool, ptr ...interface{}) error {
	if len(ptr) <= 0 {
		return errors.New("slice is empty")
	}
	sql := genSql(ptr[0], tableName)
	stm, err := Mysql.Prepare(sql)
	if err != nil {
		logrus.Errorf("执行sql：%s, 预处理失败：%v", sql, err)
		return err
	}
	defer stm.Close()
	for i := 0; i < len(ptr); i++ {
		_, err := stm.Exec(getStructValues(ptr[i])...)
		if err != nil {
			if ignoreDuplicate && strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
				logrus.Infof("insert transaction data duplicat:%v", err)
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func genSql(ptr interface{}, tableName string) string {
	v := reflect.ValueOf(ptr).Elem()
	var fieldsName []string
	var placeHolders []string
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		label := fieldInfo.Tag.Get(TAG_ORM_NAME)
		if label != "" {
			fieldsName = append(fieldsName, label)
			placeHolders = append(placeHolders, "?")
		}
	}
	return fmt.Sprintf("insert into %s (%s) values (%s)", tableName, strings.Join(fieldsName, ","), strings.Join(placeHolders, ","))
}

func getStructValues(ptr interface{}) []interface{} {
	v := reflect.ValueOf(ptr).Elem()
	var values []interface{}
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		label := fieldInfo.Tag.Get(TAG_ORM_NAME)
		if label != "" {
			values = append(values, getValue(v.Field(i)))
		}
	}
	return values
}

func getValue(val reflect.Value) interface{} {
	result := val.Interface()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = strconv.FormatInt(val.Int(), 10)
	case reflect.Struct:
		//	todo throw exception

	}
	return result

}
