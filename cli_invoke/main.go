package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/lpxxn/grpc_demo/protos"
	"github.com/lpxxn/grpc_demo/protos/model"
	"google.golang.org/grpc"
)

var port int

func init() {
	flag.IntVar(&port, "port", 10001, "grpc port")
}

func main() {
	rand.Seed(time.Now().Unix())
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	student := &model.Student{
		Id:   rand.Int63(),
		Name: randomdata.FullName(randomdata.RandomGender) + randomdata.City(),
		Age:  rand.Int31n(30),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	out := new(protos.Result)
	err = conn.Invoke(ctx, "/api.StudentSrv/NewStudent", student, out)
	if err != nil {
		panic(err)
	}
	fmt.Println("add student ", out.Code)
}
