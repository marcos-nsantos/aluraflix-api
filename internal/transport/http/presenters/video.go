package presenters

import "github.com/marcos-nsantos/aluraflix-api/internal/entity"

type Video struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func VideoResponse(video *entity.Video) *Video {
	return &Video{
		ID:          video.ID,
		Title:       video.Title,
		Description: video.Description,
		URL:         video.URL,
		CreatedAt:   video.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   video.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func VideosResponse(videos []*entity.Video) []*Video {
	var videosResponse []*Video
	for _, video := range videos {
		videosResponse = append(videosResponse, VideoResponse(video))
	}
	return videosResponse
}
