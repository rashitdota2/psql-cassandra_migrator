package main

import (
	"context"
	"log"
)

const readQuery = `select id, to_user_id, type, from_user_id, created_at 
				   from push_notifications 
				   where id between $1 and $2 order by id desc limit $3`

func (w *worker) runReader(ctx context.Context, idRange *IdRange, readerNumber int) {
	for {
		var id int

		rows, err := w.psql.Query(ctx, readQuery, idRange.Min, idRange.Max, w.readerLimit)
		if err != nil {
			log.Println("ERROR ON READER #", readerNumber, ": 'error': ", err.Error(), "; 'idRange': ", idRange)
			continue
		}

		for rows.Next() {
			var n PushNotification
			if err = rows.Scan(&id, &n.UserId, &n.Type, &n.FromUserId, &n.CreatedAt); err != nil {
				log.Println("ERROR ON READER-SCANNER #", readerNumber, ": 'error': ", err.Error(), "; 'LastScannedId': ", id)
				continue
			}

			w.writeChan <- &n
		}
		rows.Close()

		if id <= idRange.Min {
			w.readersDoneChan <- DoneSignal{}
			log.Println("READER #", readerNumber, " DONE JOB!")
			break
		}

		idRange.Max = id - 1
	}
}
