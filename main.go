package main

import (
	"context"
	"database/sql"
	"fmt"
	"likesApi/api"
	db "likesApi/db/sqlc"
	"likesApi/util"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

var config util.Config

func databaseQueryTask() {
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

    store := db.NewStore(conn)
	fmt.Println("Notification process begins")
    
    users, err := store.SelectUsersToNotify(context.Background())
    if err != nil {
        log.Fatal("Error from database" , err)
        return
    }

    for _ , usr := range users {
        log.Printf("Notified user: %d at %s" , usr , time.Now().String())
        store.CreateNotification(context.Background() , usr)
    }
}

func main() {
        var err error
	config, err = util.LoadConfig(".")


	if err != nil {
		log.Fatal("cannot load configuration:", err)
		return
	}

	c := cron.New()

	_, err = c.AddFunc("* * * * *", databaseQueryTask)

	if err != nil {
		fmt.Println("Error scheduling cron job", err)
		return
	}

	c.Start()

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

    store := db.NewStore(conn)
    gin.SetMode(gin.ReleaseMode)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err.Error())
	}
}
