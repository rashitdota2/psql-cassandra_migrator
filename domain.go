package main

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v4/pgxpool"
)

type workCfg struct {
	writeChanBufSize    int
	errChanBufSize      int
	readersCount        int
	readersIdRange      *IdRange
	writersCount        int
	errHandlersCount    int
	readersRequestLimit int
}

type worker struct {
	psql                *pgxpool.Pool
	cx                  *gocql.Session
	writeChan           chan *PushNotification
	errChan             chan *PushNotification
	readerLimit         int
	readersDoneChan     chan DoneSignal
	writersDoneChan     chan DoneSignal
	errHandlersDoneChan chan DoneSignal
}

type IdRange struct {
	Max int
	Min int
}

type PushNotification struct {
	UserId     uint64
	Type       string
	FromUserId uint64
	CreatedAt  time.Time
}
