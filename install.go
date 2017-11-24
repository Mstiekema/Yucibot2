package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/spf13/viper"
)

func main() {
  viper.SetConfigFile("./config.toml")
  viper.ReadInConfig()
  username := viper.GetString("mysql.username")
  password := viper.GetString("mysql.password")
  database := viper.GetString("mysql.database")
  
  db, _ := sql.Open("mysql", username+":"+password+"@/")
  
  db.Exec("CREATE DATABASE IF NOT EXISTS "+database+" DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci'")
  db.Exec("USE "+database)
  
  db.Exec(
    "CREATE TABLE user (" +
    "name VARCHAR(30) PRIMARY KEY," +
    "userId VARCHAR(36)," +
    "accToken VARCHAR(30)," +
    "points INT," +
    "xp INT DEFAULT 0," +
    "num_lines INT," +
    "level INT," +
    "timeOnline INT DEFAULT 0," +
    "timeOffline INT DEFAULT 0," +
    "profile_pic VARCHAR(200)," +
    "pickP BOOL DEFAULT TRUE," +
    "isMod BOOL DEFAULT FALSE," +
    "UNIQUE (name))",
  )
  
  db.Exec(
    "CREATE TABLE userstats (" +
    "userId INT PRIMARY KEY," +
    "slotLoss INT DEFAULT 0," +
    "slotWin INT DEFAULT 0," +
    "slotProfit INT DEFAULT 0," +
    "dungWin INT DEFAULT 0," +
    "dungProfit INT DEFAULT 0," +
    "roulLoss INT DEFAULT 0," +
    "roulWin INT DEFAULT 0," +
    "roulProfit INT DEFAULT 0)",
  )
  
  db.Exec(
    "CREATE TABLE emotestats (" +
    "id VARCHAR(50) PRIMARY KEY," +
    "name VARCHAR(30)," +
    "type VARCHAR(30)," +
    "uses INT DEFAULT 1)",
  )
  
  db.Exec(
    "CREATE TABLE quotes (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "name VARCHAR(30)," +
    "quote VARCHAR(500))",
  )
  
  db.Exec(
    "CREATE TABLE clr (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "name VARCHAR(30)," +
    "url VARCHAR(2083)," +
    "type VARCHAR(30))",
  )
  
  db.Exec(
    "CREATE TABLE commands (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "commName VARCHAR(50)," +
    "response VARCHAR(500)," +
    "commDesc VARCHAR(500)," +
    "commUse TEXT," +
    "level INT DEFAULT 100," +
    "points INT DEFAULT 0," +
    "cdType VARCHAR(50)," +
    "cd INT)",
  )
  
  db.Exec(
    "CREATE TABLE chatlogs (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
    "userId INT," +
    "log VARCHAR(500))",
  )
  
  db.Exec(
    "CREATE TABLE adminlogs (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
    "type VARCHAR(10)," +
    "log VARCHAR(500))",
  )
  
  db.Exec(
    "CREATE TABLE songrequest (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
    "name VARCHAR(30)," +
    "title VARCHAR(100)," +
    "thumb VARCHAR(250)," +
    "length INT," +
    "songid VARCHAR(30)," +
    "playState INT DEFAULT FALSE)",
  )
  
  db.Exec(
    "CREATE TABLE timeout (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "word VARCHAR(30)," +
    "type VARCHAR(30))",
  )
  
  db.Exec(
    "CREATE TABLE dungeon (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "user VARCHAR(30)," +
    "win BOOL DEFAULT FALSE)",
  )
  
  db.Exec(
    "CREATE TABLE module (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "moduleName VARCHAR(30)," +
    "shortName VARCHAR(30)," +
    "moduleDescription VARCHAR(500)," +
    "type VARCHAR(30)," +
    "state BOOL)",
  )
  
  db.Exec(
    "CREATE TABLE modulesettings (" +
    "id INT AUTO_INCREMENT PRIMARY KEY," +
    "moduleType VARCHAR(30)," +
    "settingName VARCHAR(30)," +
    "shortName VARCHAR(100)," +
    "value INT," +
    "message VARCHAR(500))",
  )
  
  db.Exec(
    "CREATE TABLE emotes (" +
    "emoteId VARCHAR(30) PRIMARY KEY," +
    "name VARCHAR(50)," +
    "type VARCHAR(10)," +
    "url VARCHAR(500))",
  )
  
  db.Exec(
    "CREATE TABLE timers (" +
    "name VARCHAR(50) PRIMARY KEY," +
    "online INT DEFAULT 0," +
    "offline INT DEFAULT 0," +
    "msg VARCHAR(500))",
  )
  
  db.Close()
}