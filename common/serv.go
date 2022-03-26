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

	"github.com/lpxxn/grpc_demo/protos"
	"github.com/lpxxn/grpc_demo/protos/api"
	"github.com/lpxxn/grpc_demo/protos/model"
)

type StudentSrv struct{ StudentList []*model.Student }

func (srv *StudentSrv) NewStudent(ctx context.Context, s *model.Student) (*protos.Result, error) {
	log.Println("new student in")
	if s != nil {
		srv.StudentList = append(srv.StudentList, s)
	}
	//time.Sleep(time.Hour)
	time.Sleep(time.Second * 20)
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
	const limit = 10
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
		fmt.Printf("get id: %d", in.Id)
		if l == 0 {
			fmt.Println("data is empty")
			return nil
		}
		if l > 0 {
			stream.Send(&api.QueryStudentResponse{StudentList: srv.StudentList[0:rand.Intn(l)]})
		}
	}
	return nil
}
