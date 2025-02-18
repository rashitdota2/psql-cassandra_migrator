package main

import (
	"context"
	"log"
)

const writeFailedDataQuery = `insert into failed_notifications (user_id, type, from_user_id, created_at)
                              values ($1, $2, $3, $4)`

func (w *worker) runErrHandler(ctx context.Context, errHandlerNumber int) {
	for {
		n := <-w.errChan
		if n == nil {
			w.errHandlersDoneChan <- DoneSignal{}
			log.Println("ErrHANDLER #", errHandlerNumber, " GOT NIL NOTIFICATION AND DONE JOB!")
			break
		}

		// retry to write
		err := w.cx.Query(writeQuery, n.UserId, n.CreatedAt.Unix(), n.Type, n.FromUserId).WithContext(ctx).Exec()

		if err == nil {
			continue
		}

		_, err = w.psql.Exec(ctx, writeFailedDataQuery, n.UserId, n.Type, n.FromUserId, n.CreatedAt)
		if err != nil {
			log.Println("ErrHANDLER #", errHandlerNumber, ": ERROR: ", err.Error(), ": LOST_NOTIFICATION: ", n)
		}
	}
}
