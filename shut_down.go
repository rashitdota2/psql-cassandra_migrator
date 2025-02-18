package main

import (
	"context"
	"log"
	"time"
)

type DoneSignal struct{}

func (wp *WorkerPool) WAIT(ctxCancel context.CancelFunc) {
	var doneReaders int
	var doneWriters int
	var doneErrHandlers int

	for doneReaders < wp.readersCount {
		<-wp.readersDoneChan
		doneReaders++
	}

	log.Println("ALL READERS DONE: count ", doneReaders)

	close(wp.writeChan)

	for doneWriters < wp.writersCount {
		<-wp.writersDoneChan
		doneWriters++
	}

	log.Println("ALL WRITERS DONE: count ", doneWriters)

	close(wp.errChan)

	for doneErrHandlers < wp.errHandlersCount {
		<-wp.errHandlersDoneChan
		doneErrHandlers++
	}

	log.Println("ALL ErrHANDLERS DONE: count ", doneErrHandlers)

	ctxCancel()

	time.Sleep(time.Second)

	if !wp.cx.Closed() {
		wp.cx.Close()
	}

	wp.psql.Close()

	close(wp.readersDoneChan)
	close(wp.writersDoneChan)
	close(wp.errHandlersDoneChan)
}
