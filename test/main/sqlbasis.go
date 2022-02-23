package main

// import (
// 	"database/sql"
// 	"fmt"
// )

// var envdbMap map[string]*sql.DB

// func GetEnvDbContext(connector config.DbConnector) *sql.DB {
// 	if envdbMap == nil {
// 		envdbMap = make(map[string]*sql.DB)
// 	}

// 	db, ok := envdbMap[connector.ID]
// 	if ok {
// 		return db
// 	} else {
// 		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", connector.Host, connector.Port, connector.UserName, connector.Password, connector.DatabaseName)
// 		db, err := sql.Open("postgres", connStr)

// 		envdbMap[connector.ID] = db
// 		return db
// 	}
// }
