package service

import (
	"strings"
	"sync"

	"github.com/1024casts/1024casts/pkg/constvar"

	"github.com/1024casts/1024casts/util"

	"github.com/1024casts/1024casts/model"
	"github.com/1024casts/1024casts/repository"

	"github.com/lexkong/log"
)

type VideoService struct {
	repo *repository.VideoRepo
}

func NewVideoService() *VideoService {
	return &VideoService{
		repository.NewVideoRepo(),
	}
}

func (srv *VideoService) CreateVideo(video model.VideoModel) (id uint64, err error) {
	if video.CoverKey != "" && !strings.HasPrefix(video.CoverKey, "/") {
		video.CoverKey = "/" + video.CoverKey
	}
	id, err = srv.repo.CreateVideo(video)

	if err != nil {
		return id, err
	}

	return id, nil
}

func (srv *VideoService) GetVideoById(id int) (*model.VideoModel, error) {
	Video, err := srv.repo.GetVideoById(id)

	if err != nil {
		return Video, err
	}

	return Video, nil
}

func (srv *VideoService) GetVideoList(courseId uint64, isGroup bool) ([]*model.VideoModel, error) {
	videos := make([]*model.VideoModel, 0)

	videos, err := srv.repo.GetVideoList(courseId, isGroup)
	if err != nil {
		log.Warnf("[video] get video list err, course_id: %d", courseId)
		return nil, err
	}

	for _, video := range videos {
		video.CoverUrl = util.GetVideoCover(video.CoverKey)
		video.Mp4URL = util.GetVideoUrl(video.Mp4Key)
		video.DurationStr = util.ResolveVideoDuration(video.Duration)
		video.PublishedAtStr = util.FormatTime(video.PublishedAt)
	}

	return videos, nil
}

func (srv *VideoService) GetVideoByCourseIdAndEpisodeId(courseId uint64, episodeId int) (*model.VideoModel, error) {
	video := new(model.VideoModel)

	video, err := srv.repo.GetVideoByCourseIdAndEpisodeId(courseId, episodeId)
	if err != nil {
		log.Warnf("[video] get video list err, course_id: %d", courseId)
		return nil, err
	}

	video.Mp4URL = util.GetVideoUrl(video.Mp4Key)
	video.DurationStr = util.ResolveVideoDuration(video.Duration)
	video.PublishedAtStr = util.FormatTime(video.PublishedAt)

	return video, nil
}

func (srv *VideoService) GetVideoListPagination(courseId uint64, name string, offset, limit int) ([]*model.VideoModel, uint64, error) {
	infos := make([]*model.VideoModel, 0)

	Videos, count, err := srv.repo.GetVideoListPagination(courseId, name, offset, limit)
	if err != nil {
		return nil, count, err
	}

	ids := []uint64{}
	for _, Video := range Videos {
		ids = append(ids, Video.Id)
	}

	wg := sync.WaitGroup{}
	VideoList := model.VideoList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.VideoModel, len(Videos)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// Improve query efficiency in parallel
	for _, c := range Videos {
		wg.Add(1)
		go func(Video *model.VideoModel) {
			defer wg.Done()

			//shortId, err := util.GenShortId()
			//if err != nil {
			//	errChan <- err
			//	return
			//}

			VideoList.Lock.Lock()
			defer VideoList.Lock.Unlock()

			Video.Mp4URL = util.GetQiNiuPrivateAccessUrl(Video.Mp4Key, constvar.MediaTypeVideo, 0, 0)

			VideoList.IdMap[Video.Id] = Video
		}(c)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	for _, id := range ids {
		infos = append(infos, VideoList.IdMap[id])
	}

	return infos, count, nil
}

// get video total count
func (srv *VideoService) GetVideoTotalCount() (int, error) {
	count, err := srv.repo.GetVideoTotalCount()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// get video total duration 总秒数
func (srv *VideoService) GetVideoTotalDuration() (int, error) {
	count, err := srv.repo.GetVideoTotalDuration()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (srv *VideoService) UpdateVideo(VideoMap map[string]interface{}, id int) error {
	err := srv.repo.UpdateVideo(VideoMap, id)

	if err != nil {
		return err
	}

	return nil
}
