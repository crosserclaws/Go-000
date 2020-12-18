package biz

import (
	"fmt"

	"week04/internal/domain"
)

// VideoRepository is a repo interface of video.
type VideoRepository interface {
	GetVideoInformation(int64) (*domain.VideoInformation, error)
}

// VideoInfoBiz contains business logics about video information.
type VideoInfoBiz struct {
	r VideoRepository
}

// NewVideoInfoBiz creates a new instance.
func NewVideoInfoBiz(r VideoRepository) *VideoInfoBiz {
	return &VideoInfoBiz{r: r}
}

// GetVideoInfo looks up a video's info by its ID.
func (v *VideoInfoBiz) GetVideoInfo(id int64) (*domain.VideoInformation, error) {
	if id < 0 {
		return nil, fmt.Errorf("Invalid video id: %v", id)
	}

	info, err := v.r.GetVideoInformation(id)
	if err != nil {
		return nil, err
	} else if info.IsReported {
		// Logic: check if a video is reported.
		return nil, fmt.Errorf("The video id=%v is reported. Cannot retrieve any information from it", id)
	} else {
		return info, nil
	}
}
