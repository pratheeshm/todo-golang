package main

import (
	"fmt"
	"net/http"

	"github.com/pratheeshm/todo-golang/task/usecase"

	"github.com/pratheeshm/todo-golang/task/repository"

	taskdeliver "github.com/pratheeshm/todo-golang/task/delivery/http"

	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("config/")
	viper.SetConfigName("app")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
}
func main() {
	db, err := mustInitDB()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Info("Connected to DB successfully")
	tr := repository.NewPostgresTaskRepository(db)
	tu := usecase.NewTaskUsecase(tr)
	h := taskdeliver.NewTaskHandler(tu)
	err = http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("server.port")), h)
	if err != nil {
		log.Panic(err)
	}
}
func mustInitDB() (*sql.DB, error) {
	host := viper.GetString("database.host")
	user := viper.GetString("database.username")
	pass := viper.GetString("database.password")
	port := viper.GetInt("database.port")
	ssl := viper.GetString("database.sslmode")
	dbname := viper.GetString("database.dbname")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, ssl)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}
