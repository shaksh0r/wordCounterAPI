package main

import (
	"fmt"
)

type Request struct{
	name string
	age int
}

func handler( r *Request, sem chan int){
	sem <- 1
	fmt.Printf("name: %v age: %v \n",r.name,r.age)
	<-sem
}

func server(queue chan *Request,sem chan int){
	for {
		req := <- queue
		go handler(req,sem)
	}
}

func main(){

	queue := make(chan *Request,10)
	sem := make(chan int,100)

	go server(queue,sem)

	for i:= 0;i<1000;i++{
		req := Request{name:"shakshor",age:i}
		queue <- &req
	}


}