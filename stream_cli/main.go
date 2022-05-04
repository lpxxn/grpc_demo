package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/lpxxn/grpc_demo/protos/api"
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
	//revAllStudent(c)
	revStudentInfo(c)
}

func revAllStudent(c api.StudentSrvClient) {
	stream, err := c.AllStudent(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}

	for {
		rev, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
		for _, item := range rev.StudentList {
			log.Printf("Got message %+v \n", item)
		}
		log.Println(strings.Repeat("--", 20))
	}
}

func revStudentInfo(c api.StudentSrvClient) {
	stream, err := c.StudentInfo(context.Background())
	if err != nil {
		panic(err)
	}
	var w sync.WaitGroup
	go func() {
		w.Add(1)
		for {
			rev, err := stream.Recv()
			if err == io.EOF {
				w.Done()
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			for _, item := range rev.StudentList {
				log.Printf("Got message %+v \n", item)
			}
			log.Println(strings.Repeat("--", 20))
		}
	}()
	send := func(id int64) {
		if err := stream.Send(&api.QueryStudent{Id: id}); err != nil {
			panic(err)
		}
	}
	for i := 0; i < 100; i++ {
		send(int64(i))
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}
	stream.CloseSend()
	w.Wait()
}
