package qiniu

import (
	"1024casts/backend/model"

	"github.com/qiniu/api.v7/storage"
)

type CreateRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	CoverImage  string `json:"cover_image"`
	UserId      int    `json:"user_id"`
	IsPublish   int    `json:"is_publish"`
}

type CreateResponse struct {
	Id uint64 `json:"id"`
}

type ListRequest struct {
	Name  string `json:"name"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type UploadResponse struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
}

type ListResponse struct {
	TotalCount uint64             `json:"totalCount"`
	List       []storage.ListItem `json:"list"`
	NextMaker  string             `json:"next_maker"`
}

type SwaggerListResponse struct {
	TotalCount uint64               `json:"totalCount"`
	List       []model.CommentModel `json:"list"`
}