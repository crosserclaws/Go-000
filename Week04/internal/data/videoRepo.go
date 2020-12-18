package data

import (
	"fmt"
	"week04/internal/biz"
	"week04/internal/domain"
)

var _ biz.VideoRepository = new(VideoRepo)

// VideoRepo implements biz.VideoRepository.
type VideoRepo struct {
	// NOTE: The storage can be a Redis, we just fake the behavior here for simplicity.
	c *VideoCountStorage
	// NOTE: The storage can be a RDB such as MySQL, we just fake the behavior here for simplicity.
	m *VideoMetaStorage
}

// NewVideoRepo creates a new instance.
func NewVideoRepo(c *VideoCountStorage, m *VideoMetaStorage) *VideoRepo {
	return &VideoRepo{c, m}
}

// GetVideoInformation implements biz.VideoRepository.
func (v *VideoRepo) GetVideoInformation(id int64) (*domain.VideoInformation, error) {
	count, err := v.c.GetCount(id)
	if err != nil {
		return nil, err
	}

	meta, err := v.m.GetMetadata(id)
	if err != nil {
		return nil, err
	}
	return &domain.VideoInformation{
		Count:      count,
		Name:       meta.Name,
		IsReported: meta.IsReported,
	}, nil
}

// VideoCountStorage contains counts of videos.
type VideoCountStorage struct {
}

// NewVideoCountStorage creates a new instance.
func NewVideoCountStorage() *VideoCountStorage {
	return &VideoCountStorage{}
}

// GetCount return a video's view count by given ID.
func (v *VideoCountStorage) GetCount(id int64) (int64, error) {
	// NOTE: The storage can be a Redis, we just fake the behavior here for simplicity.
	return id, nil
}

// VideoMetaStorage contains static meta data of videos.
type VideoMetaStorage struct {
}

// NewVideoMetaStorage creates a new instance.
func NewVideoMetaStorage() *VideoMetaStorage {
	return &VideoMetaStorage{}
}

// GetMetadata return a video's meta data which should not be changed frequently.
func (v *VideoMetaStorage) GetMetadata(id int64) (*VideoMeta, error) {
	// NOTE: The storage can be a RDB such as MySQL, we just fake the behavior here for simplicity.
	return &VideoMeta{
		Name:       fmt.Sprintf("Video-%v", id),
		IsReported: (id & 3) == 0,
	}, nil
}

// VideoMeta is a PO(persistent object).
type VideoMeta struct {
	Name       string
	IsReported bool
}
