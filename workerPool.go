package main

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WorkerPool struct {
	psql                *pgxpool.Pool
	cx                  *gocql.Session
	writeChan           chan *PushNotification
	errChan             chan *PushNotification
	readersCount        int
	readersIdRanges     map[int]*IdRange
	writersCount        int
	errHandlersCount    int
	readerLimit         int
	readersDoneChan     chan DoneSignal
	writersDoneChan     chan DoneSignal
	errHandlersDoneChan chan DoneSignal
}

func NewWorkerPool(cfg *workCfg, psql *pgxpool.Pool, cx *gocql.Session) *WorkerPool {
	return &WorkerPool{
		psql:                psql,
		cx:                  cx,
		writeChan:           make(chan *PushNotification, cfg.writeChanBufSize),
		errChan:             make(chan *PushNotification, cfg.errChanBufSize),
		readersCount:        cfg.readersCount,
		readersIdRanges:     calculateIdRanges(cfg.readersIdRange, cfg.readersCount),
		writersCount:        cfg.writersCount,
		errHandlersCount:    cfg.errHandlersCount,
		readerLimit:         cfg.readersRequestLimit,
		readersDoneChan:     make(chan DoneSignal),
		writersDoneChan:     make(chan DoneSignal),
		errHandlersDoneChan: make(chan DoneSignal),
	}
}

func (wp *WorkerPool) RUN(ctx context.Context) {
	w := &worker{
		psql:                wp.psql,
		cx:                  wp.cx,
		writeChan:           wp.writeChan,
		errChan:             wp.errChan,
		readerLimit:         wp.readerLimit,
		readersDoneChan:     wp.readersDoneChan,
		writersDoneChan:     wp.writersDoneChan,
		errHandlersDoneChan: wp.errHandlersDoneChan,
	}

	for i := 1; i <= wp.readersCount; i++ {
		go w.runReader(ctx, wp.readersIdRanges[i], i)
	}

	for i := 1; i <= wp.writersCount; i++ {
		go w.runWriter(ctx, i)
	}

	for i := 1; i <= wp.errHandlersCount; i++ {
		go w.runErrHandler(ctx, i)
	}
}
