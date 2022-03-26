package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
	"google.golang.org/grpc"

	"github.com/lpxxn/grpc_demo/protos/api"
	"github.com/lpxxn/grpc_demo/protos/model"
)

var port int

func init() {
	flag.IntVar(&port, "port", 10001, "grpc port")
}

func main() {
	rand.Seed(time.Now().Unix())
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	//conn, err := grpc.Dial(fmt.Sprintf("192.168.10.94:%d", 62937), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewStudentSrvClient(conn)

	for i := 0; i < 1; i++ {
		go func(idx int) {
			fmt.Println("new student", idx)
			r, err := c.NewStudent(context.Background(), &model.Student{
				Id:   rand.Int63(),
				Name: randomdata.FullName(randomdata.RandomGender),
				Age:  rand.Int31n(30),
			})
			if err != nil {
				panic(err)
			}
			fmt.Println("add student ", r.Code)
			//time.Sleep(time.Second * 2)
		}(i)
	}
	time.Sleep(time.Minute * 10)
}
