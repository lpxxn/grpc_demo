package common

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/lpxxn/grpc_demo/protos"
	"github.com/lpxxn/grpc_demo/protos/api"
	"github.com/lpxxn/grpc_demo/protos/model"
)

type StudentSrv struct {
	StudentList []*model.Student
	Version     string
}

func (srv *StudentSrv) NewStudent(ctx context.Context, s *model.Student) (*protos.Result, error) {
	log.Println("new student in", srv.Version)
	{
		if meta, ok := metadata.FromIncomingContext(ctx); ok {
			log.Println(meta)
			log.Println("user:", meta.Get("user"))
		}
		header := metadata.Pairs("header", "headerVal")
		grpc.SendHeader(ctx, header)
		trailer := metadata.Pairs("trailer", "trailerVal")
		grpc.SetTrailer(ctx, trailer)
	}
	if s != nil {
		srv.StudentList = append(srv.StudentList, s)
	}
	if err := checkCtx(ctx); err != nil {
		return &protos.Result{
			Code: "false",
			Desc: err.Error(),
		}, nil
	}
	return &protos.Result{
		Code: "OK",
		Desc: randomdata.FullName(randomdata.RandomGender) + randomdata.Address(),
	}, nil
}

func checkCtx(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Println("ctx canceled")
		return ctx.Err()
	default:
		return nil
	}
}

func (srv *StudentSrv) StudentByID(context.Context, *api.QueryStudent) (*api.QueryStudentResponse, error) {
	l := len(srv.StudentList)
	rev := &api.QueryStudentResponse{StudentList: srv.StudentList}
	srv.StudentList = srv.StudentList[l:]
	return rev, nil
}

func (srv *StudentSrv) AllStudent(e *empty.Empty, rev api.StudentSrv_AllStudentServer) error {
	const limit = 2
	data := &api.QueryStudentResponse{}
	curr := srv.StudentList
	for _, item := range curr {
		data.StudentList = append(data.StudentList, item)
		if len(data.StudentList) == limit {
			if err := rev.Send(data); err != nil {
				fmt.Printf("send error %#v", err)
			}
			data.StudentList = data.StudentList[:0]
		}
	}
	if len(data.StudentList) > 0 {
		if err := rev.Send(data); err != nil {
			fmt.Printf("send error %#v", err)
		}
	}
	return nil
}
func (srv *StudentSrv) StudentInfo(stream api.StudentSrv_StudentInfoServer) error {
	l := len(srv.StudentList)
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		fmt.Printf("get id: %d\n", in.Id)
		if l == 0 {
			fmt.Println("data is empty")
			return nil
		}
		if l > 0 {
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			stream.Send(&api.QueryStudentResponse{StudentList: srv.StudentList[0:rand.Intn(l)]})
		}
	}
	return nil
}

func (srv *StudentSrv) QueryStudents(student *api.QueryStudent, server api.StudentSrv_QueryStudentsServer) error {
	//TODO implement me
	panic("implement me")
}
