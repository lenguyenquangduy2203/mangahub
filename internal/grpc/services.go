package grpc

import (
	"context"
	"errors"
	"mangahub/internal/api-server/mangas"
	"mangahub/internal/api-server/users"
	"mangahub/pkg/models"
	mangapb "mangahub/proto/manga"
	"strings"

	"mangahub/internal/tcp/client"
	update "mangahub/internal/tcp/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MangaServiceServer struct {
	mangapb.UnimplementedMangaServiceServer
	userService           users.UserLibraryService
	mangaService          mangas.MangaService
	tcpProgressSyncClient client.ProgressSyncClient
}

func NewMangaServiceGrpcServer(userService users.UserLibraryService, mangaService mangas.MangaService, tcpClient client.ProgressSyncClient) *MangaServiceServer {
	return &MangaServiceServer{
		userService:           userService,
		mangaService:          mangaService,
		tcpProgressSyncClient: tcpClient,
	}
}

func (s *MangaServiceServer) GetManga(ctx context.Context, req *mangapb.GetMangaRequest) (*mangapb.MangaResponse, error) {
	manga, err := s.mangaService.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, mangas.ErrMangaNotExistInDatabase) {
			return nil, status.Errorf(codes.NotFound, "manga not found: %s", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "failed to get manga: %v", err)
	}

	return &mangapb.MangaResponse{
		Id:            manga.ID,
		Title:         manga.Title,
		Author:        manga.Author,
		Genres:        manga.Genres,
		Status:        manga.Status,
		TotalChapters: int32(manga.TotalChapters),
		Description:   manga.Description,
	}, nil
}

func (s *MangaServiceServer) SearchManga(ctx context.Context, req *mangapb.SearchRequest) (*mangapb.SearchResponse, error) {
	query := models.MangaSearchQuery{
		Title:  strings.TrimSpace(req.GetTitle()),
		Author: strings.TrimSpace(req.GetAuthor()),
		Genre:  strings.TrimSpace(req.GetGenre()),
		Status: strings.ToUpper(strings.TrimSpace(req.GetStatus())),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}

	result, err := s.mangaService.Find(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search failed: %v", err)
	}

	mangaItems := make([]*mangapb.MangaListItem, 0, len(result.Results))
	for _, item := range result.Results {
		mangaItems = append(mangaItems, &mangapb.MangaListItem{
			Id:            item.ID,
			Title:         item.Title,
			TotalChapters: int32(item.TotalChapters),
			Status:        item.Status,
		})
	}

	return &mangapb.SearchResponse{
		Total:   result.Total,
		Limit:   int32(result.Limit),
		Offset:  int32(result.Offset),
		Results: mangaItems,
	}, nil
}

func (s *MangaServiceServer) UpdateProgress(ctx context.Context, req *mangapb.ProgressRequest) (*mangapb.ProgressResponse, error) {
	err := s.userService.UpdateUserReadingProgress(ctx, req.GetUserId(), req.GetMangaId(), int(req.GetChapter()))
	if err != nil {
		if errors.Is(err, mangas.ErrMangaNotExistInDatabase) {
			return nil, status.Error(codes.NotFound, "manga not found in database")
		}
		if errors.Is(err, users.ErrMangaNotExistInUserLibrary) {
			return nil, status.Error(codes.NotFound, "manga not found in user library")
		}
		if errors.Is(err, mangas.ErrInValidCurrentChapter) {
			return nil, status.Error(codes.InvalidArgument, "invalid chapter number")
		}
		return nil, status.Errorf(codes.Internal, "update progress failed: %v", err)
	}

	s.tcpProgressSyncClient.Send(update.ProgressUpdate{
		UserID:  req.GetUserId(),
		MangaID: req.GetMangaId(),
		Chapter: int(req.GetChapter()),
	})

	return &mangapb.ProgressResponse{
		Ok: true,
	}, nil
}
