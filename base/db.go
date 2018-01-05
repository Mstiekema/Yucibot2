package base

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/spf13/viper"
)

func Conn() (db *sql.DB) {
  viper.SetConfigFile("./config.toml")
  err := viper.ReadInConfig()
  if err != nil {
    fmt.Println(err)
  }
  username := viper.GetString("mysql.username")
  password := viper.GetString("mysql.password")
  table := viper.GetString("mysql.database")
  
  db, err = sql.Open("mysql", username+":"+password+"@/"+table)
  if err != nil {
		panic(err.Error())
	}
  return
}

func Query(qry string) (result string) {
  var db = Conn()
  stmtOut, err := db.Prepare(qry)
  if err != nil {
		panic(err.Error())
	}
  stmtOut.QueryRow().Scan(&result)
  db.Close()
  return
}

func Update(tbl, row, name, where string) {
  var db = Conn()
  db.Exec("UPDATE "+tbl+" SET "+row+" WHERE "+name+" = "+where)
  db.Close()
}

func Delete(tbl, row, where string) {
  var db = Conn()
  db.Exec("DELETE FROM "+tbl+" WHERE "+row+" = '"+where+"'")
  db.Close()
}

func Insert(tbl, set string) {
  var db = Conn()
  db.Exec("INSERT INTO "+tbl+" VALUES "+set)
  db.Close()
}