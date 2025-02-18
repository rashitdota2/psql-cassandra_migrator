package main

import (
	"context"
	"log"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v4/pgxpool"
)

// example of cfg
var workConfig = &workCfg{
	readersCount: 3,
	readersIdRange: &IdRange{
		Min: 1,
		Max: 430,
	},
	readersRequestLimit: 3,

	writersCount:     2,
	writeChanBufSize: 10,

	errHandlersCount: 2,
	errChanBufSize:   10,
}

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())

	// connect to psql and cassandra with your cfg, settings, etc.
	var psql *pgxpool.Pool
	var cqlClient *gocql.Session

	log.Println("ALL CONNECTIONS ESTABLISHED. START JOBS!!!")

	workerPool := NewWorkerPool(workConfig, psql, cqlClient)

	workerPool.RUN(ctx)

	workerPool.WAIT(ctxCancel)

	log.Println("ALL JOBS DONE!!!")
}
