package grpc

import (
	"mangahub/internal/mangas"
	"mangahub/internal/users"
	mangapb "mangahub/proto/manga"
)

type MangaServiceServer struct {
	mangapb.UnimplementedMangaServiceServer
	userRepo  users.Repository
	mangaRepo mangas.Repository
}

func NewMangaServiceGrpcServer(userRepo users.Repository, mangaRepo mangas.Repository) *MangaServiceServer {
	return &MangaServiceServer{
		userRepo:  userRepo,
		mangaRepo: mangaRepo,
	}
}
