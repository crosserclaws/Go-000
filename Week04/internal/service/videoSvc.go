package service

// Implement gRPC interface.
// Call biz

import (
	"context"
	"log"

	pb "week04/api/video/v1"
	"week04/internal/biz"
	"week04/internal/domain"
)

// VideoInfoService contains video information.
type VideoInfoService struct {
	pb.UnimplementedVideoInformerServer
	b *biz.VideoInfoBiz
}

// NewVideoInfoService creates a new instance.
func NewVideoInfoService(b *biz.VideoInfoBiz) *VideoInfoService {
	return &VideoInfoService{b: b}
}

// GetVideoInfo implements VideoInformer.GetVideoInfo
func (v *VideoInfoService) GetVideoInfo(ctx context.Context, in *pb.VideoInfoRequest) (*pb.VideoInfoReply, error) {
	log.Println("[Video info svc] Get video info: ID=", in.Id)
	r, err := v.b.GetVideoInfo(in.Id)
	return VideoInformationToReply(r), err
}

// VideoInformationToReply converts the DO to DTO
func VideoInformationToReply(info *domain.VideoInformation) *pb.VideoInfoReply {
	if info == nil {
		return nil
	}
	return &pb.VideoInfoReply{Name: info.Name, Count: info.Count}
}
