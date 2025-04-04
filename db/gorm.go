package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func F_连接MYSQL(drive string) *gorm.DB {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = dir + "/db.json"
	buf, _ := ioutil.ReadFile(dir)
	s := gjson.Get(string(buf), drive)
	server := s.Get("server").String()
	port := s.Get("port").Int()
	user := s.Get("user").String()
	password := s.Get("password").String()
	database := s.Get("database").String()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, server, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println(dsn)
		panic(err)
	}
	return db
}

func F_连接SqlServer(drive string) (string, *gorm.DB) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = dir + "/db.json"
	buf, _ := ioutil.ReadFile(dir)
	s := gjson.Get(string(buf), drive)
	fmt.Println(string(buf))
	server := s.Get("server").String()
	port := s.Get("port").Int()
	user := s.Get("user").String()
	password := s.Get("password").String()
	database := s.Get("database").String()
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, server, strconv.Itoa(int(port)), database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return dsn, db
}
