package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/lpxxn/grpc_demo/protos/api"
	"github.com/lpxxn/grpc_demo/protos/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	c := api.NewStudentSrvClient(conn)
	student := &model.Student{
		Id:    rand.Int63(),
		Value: randomdata.FullName(randomdata.RandomGender) + randomdata.City(),
		Age:   rand.Int31n(30),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	md := metadata.New(map[string]string{
		"user": "abc",
	})
	// 使用Pairs
	md2 := metadata.Pairs("k", "v", "k2", "v2")
	md = metadata.Join(md, md2)
	ctx = metadata.NewOutgoingContext(ctx, md)
	ctx = metadata.AppendToOutgoingContext(ctx, "k3", "v3")

	revHeader := metadata.MD{}
	revTrailer := metadata.MD{}
	r, err := c.NewStudent(ctx, student, grpc.Header(&revHeader), grpc.Trailer(&revTrailer))

	if err != nil {
		panic(err)
	}
	fmt.Println("add student ", r.Code, revHeader, revTrailer)
}
