package base

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/spf13/viper"
)

// Using the old database for now. Also only able to pull one thing at the time. Will improve in the future
func Query(tbl, row, where string) (result string) {
  viper.SetConfigFile("./config.toml")
  err := viper.ReadInConfig()
  if err != nil {
    fmt.Println(err)
  }
  username := viper.GetString("mysql.username")
  password := viper.GetString("mysql.password")
  table := viper.GetString("mysql.database")
  
  db, err := sql.Open("mysql", username+":"+password+"@/"+table)
  if err != nil {
		panic(err.Error())
	}
  stmtOut, err := db.Prepare("SELECT "+row+" FROM "+tbl+" WHERE name = ?")
  if err != nil {
		panic(err.Error())
	}
  stmtOut.QueryRow(where).Scan(&result) 
  return
}