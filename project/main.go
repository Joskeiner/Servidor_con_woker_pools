package main

import (
	"fmt"
	"time"
)

type Job struct {
	Name   string        // nombre de trabajo
	Delay  time.Duration // la duracion que tendra el durmiendo
	Number int           // numero que usara la funcion de fibonacci
}
type Worker struct {
	Id         int           // sera el numero de identificacion de el canal
	JobQueue   chan Job      // este canal es  que estaremos utilizando para leer los jobs
	WorkerPool chan chan Job // resivira nacales jobs
	QuitChan   chan bool     // este canal se usara  para cerrar los wokers
}

// Constructor para Worker
func NewWorker(id int, wokerPool chan chan Job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan Job),
		WorkerPool: wokerPool,
		QuitChan:   make(chan bool),
	}
}

// Metodos para leer los canales y hacer multiplexacion  de canales
func (w Worker) Start() {
	go func() {
		for { // lee de  w.JobQueue
			w.WorkerPool <- w.JobQueue // envio de chan Job
			select {
			case job := <-w.JobQueue:
				fmt.Printf("Worker with id %d Started \n", w.Id)
				fib := Fibonacci(job.Number)
				time.Sleep(job.Delay)
				fmt.Printf("Worker with id %d Finished with result %d\n", w.Id, fib)
			case <-w.QuitChan: // cerrar el canal Worker
				fmt.Printf("Worker with id %d Stopped\n", w.Id)
			}
		}
	}()
}

// cerrar el canal Worker
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int      // sera la cantidad de buffer para los canales
	JobQueue   chan Job // cola de trabajo
}

// construtor Dispatcher
func NewDispatcher(jobQueue chan Job, maxWorker int) *Dispatcher {
	woker := make(chan chan Job, maxWorker)
	return &Dispatcher{
		JobQueue:   jobQueue,
		MaxWorkers: maxWorker,
		WorkerPool: woker,
	}
}

// metodo para encolar los job para que los lea en funcion Strat
func (d *Dispatcher) Dispatcher() {
	for {
		select {
		case job := <-d.JobQueue:
			go func() {
				workerJobQueue := <-d.WorkerPool
				workerJobQueue <- job
			}()
		}
	}
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}
