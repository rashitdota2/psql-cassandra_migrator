package main

import (
	"context"
	"log"
)

const writeQuery = `insert into push_notifications (user_id, created_at, type, from_user_id)
					values (?, ?, ?, ?)`

func (w *worker) runWriter(ctx context.Context, writerNumber int) {
	for {
		n := <-w.writeChan
		if n == nil {
			w.writersDoneChan <- DoneSignal{}
			log.Println("WRITER #", writerNumber, " GOT NIL NOTIFICATION AND DONE JOB!")
			break
		}

		err := w.cx.Query(writeQuery, n.UserId, n.CreatedAt.UnixMilli(), n.Type, n.FromUserId).WithContext(ctx).Exec()

		if err != nil {
			w.errChan <- n
		}
	}
}
