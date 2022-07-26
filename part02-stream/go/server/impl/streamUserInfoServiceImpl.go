package impl

import (
	"context"
	"errors"
	"fmt"
	"grpc-in-action/part02-stream/go/server/pb"
	"io"
	"log"
)

type StreamUserInfoServiceImpl struct {
	UserLocationTmp map[string][]UserLocTmp
	UserInfoData    map[string]*UserInfo
	pb.UnimplementedUserInfoServiceServer
}

type UserLocTmp struct {
	Loc         string
	Temperature string
}

type UserInfo struct {
	Id   string
	Name string
	Age  int32
}

func (server *StreamUserInfoServiceImpl) GetUserInfo(ctx context.Context, in *pb.GetUserInfoReq) (*pb.GetUserInfoRes, error) {

	log.Printf("go stream server收到请求")
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

func (server *StreamUserInfoServiceImpl) ReportUserBehavior(stream pb.UserInfoService_ReportUserBehaviorServer) error {

	res := pb.ReportUserBehaviorRes{
		RetCode:       0,
		RetMsg:        "收到收到~",
		ReceivedCount: 0,
	}
	for {
		behaviorData, err := stream.Recv()
		if err == io.EOF {
			// 结束符 客户端停止发送
			return stream.SendAndClose(&res)
		}
		if err != nil {
			return err
		}
		//此处有并发问题 只是简单演示
		res.ReceivedCount += 1
		log.Printf("用户 ID : %s -行为： %s", behaviorData.Id, behaviorData.Behavior)
	}

}

func (server *StreamUserInfoServiceImpl) SupplyUserChangeInfo(req *pb.SupplyUserChangeInfoReq,
	res pb.UserInfoService_SupplyUserChangeInfoServer) error {
	//获取用户的数据
	if server.UserLocationTmp == nil {
		server.initUserLocationTmp()
	}

	id := req.Id
	if len(id) == 0 {
		return errors.New("id不能为空")
	}
	//获取该用户的历史移动轨迹
	locTmps := server.UserLocationTmp[id]
	for i, tmp := range locTmps {
		log.Printf("i:%d,data:%v", i, tmp)
		err := res.Send(&pb.SupplyUserChangeInfoRes{
			Id:          id,
			Loc:         tmp.Loc,
			Temperature: tmp.Temperature,
		})
		if err != nil {
			return fmt.Errorf("error sending message to stream : %v", err)
		}
	}
	return nil
}

func (server *StreamUserInfoServiceImpl) ExchangeUserInfo(ex pb.UserInfoService_ExchangeUserInfoServer) error {
	for {
		recv, err := ex.Recv()
		if err == io.EOF {
			err := ex.Send(&pb.ExchangeUserInfoRes{
				Id:  "结束汇报",
				Loc: "结束汇报",
			})
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("用户Id:%s, lat:%v,lng:%v", recv.Id, recv.Lat, recv.Lng)
		err = ex.Send(&pb.ExchangeUserInfoRes{
			Id:  recv.Id,
			Loc: "收到，您当前的位置是北京",
		})
		if err != nil {
			return err
		}
	}
}

func (server *StreamUserInfoServiceImpl) initUserInfoData() {
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

func (server *StreamUserInfoServiceImpl) initUserLocationTmp() {
	server.UserLocationTmp = make(map[string][]UserLocTmp)
	var data []UserLocTmp
	for i := 0; i < 10; i++ {
		t := UserLocTmp{
			Loc:         "北京",
			Temperature: "26.5",
		}
		data = append(data, t)
	}
	server.UserLocationTmp["1"] = data
	server.UserLocationTmp["2"] = data
	server.UserLocationTmp["3"] = data
}
