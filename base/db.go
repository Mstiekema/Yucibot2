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

func Query(tbl, row, where string) (result string) {
  var db = Conn()
  stmtOut, err := db.Prepare("SELECT "+row+" FROM "+tbl+" WHERE name = ?")
  if err != nil {
		panic(err.Error())
	}
  stmtOut.QueryRow(where).Scan(&result)
  return
}

func Update(tbl, row, where string) {
  var db = Conn()
  db.Exec("UPDATE "+tbl+" SET "+row+" WHERE name = "+where)
}