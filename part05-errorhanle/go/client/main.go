package main

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpc-in-action/part05-errorhanle/go/client/pb"
	"log"
)

func main() {

	dial, err := grpc.Dial("localhost:10085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewBlendServiceClient(dial)
	blend, err := client.Blend(context.Background(), &pb.BlendReq{
		UserId: "-1",
	})
	if err != nil {
		log.Printf("err is :%+v",err)
		errorCode := status.Code(err)
		if errorCode == codes.InvalidArgument {
			log.Printf("Invalid Argument Error : %s", errorCode)
			errorStatus := status.Convert(err)
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				case *errdetails.BadRequest_FieldViolation:
					log.Printf("Request Field Invalid: %s", info)
				default:
					log.Printf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Printf("Unhandled error : %s ", errorCode)
		}
	}else {
		log.Printf("blend:%v",blend.Nid)
	}

}
