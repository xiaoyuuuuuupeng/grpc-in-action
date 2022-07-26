package serviceimpl

import (
	"context"
	"errors"
	"grpc-in-action/part1/go/server/pb"
	"log"
)

type UserInfo struct {
	Id   string
	Name string
	Age  int32
}

type UserInfoServiceImpl struct {
	UserInfoData map[string]*UserInfo
	pb.UnimplementedUserInfoServiceServer
}

func (server *UserInfoServiceImpl) GetUserInfo(ctx context.Context, in *pb.GetUserInfoReq) (*pb.GetUserInfoRes, error) {

	log.Printf("go server收到请求")
	//初始化一个map
	if server.UserInfoData == nil {
		server.initUserInfoData()
	}
	data := server.UserInfoData[in.Id]
	if data == nil {
		return nil, errors.New("该id不存在")
	}
	res := &pb.GetUserInfoRes{
		Id:   data.Id,
		Name: data.Name,
		Age:  data.Age,
	}

	return res, nil
}

func (server *UserInfoServiceImpl) initUserInfoData() {
	server.UserInfoData = make(map[string]*UserInfo)

	server.UserInfoData["1"] = &UserInfo{
		Id:   "1",
		Name: "小美",
		Age:  18,
	}
	server.UserInfoData["2"] = &UserInfo{
		Id:   "2",
		Name: "小刚",
		Age:  28,
	}
	server.UserInfoData["3"] = &UserInfo{
		Id:   "3",
		Name: "小王",
		Age:  20,
	}
}
