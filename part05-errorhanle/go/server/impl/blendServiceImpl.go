package impl

import (
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-in-action/part05-errorhanle/go/server/pb"
	"log"
)

type BlendService struct {
	pb.UnimplementedBlendServiceServer
}

func (*BlendService) Blend(context context.Context, req *pb.BlendReq) (*pb.BlendRes, error)  {

	if req.UserId == "-1" {
		log.Printf("UserId is invalid! -> Received User ID %s", req.UserId)
		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		ds, err := errorStatus.WithDetails(
			&errdetails.BadRequest_FieldViolation{
				Field:"ID",
				Description: fmt.Sprintf("UserId received is not valid :%s ", req.UserId),
			},
		)
		if err != nil {
			return nil, errorStatus.Err()
		}
		return nil, ds.Err()
	}else {
		nids :=[]string{"1","2","3"}

		return &pb.BlendRes{
			Nid:nids,
		},nil
	}
}
