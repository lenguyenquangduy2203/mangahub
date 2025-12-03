package tests

import (
	"context"
	"net"
	"testing"

	mygrpc "mangahub/internal/grpc"
	"mangahub/internal/mangas"
	update "mangahub/internal/tcp/models"
	"mangahub/internal/users"
	"mangahub/pkg/models"
	mangapb "mangahub/proto/manga"

	"mangahub/internal/tcp/client"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

//
// ────────────────────────────────────────────────
//   TEST DOUBLE DEFINITIONS
// ────────────────────────────────────────────────
//

// ---- Manga service stub ----
type MangaServiceStub struct {
	mangas.MangaService

	GetFn  func(ctx context.Context, id string) (models.Manga, error)
	FindFn func(ctx context.Context, q models.MangaSearchQuery) (models.PaginatedMangaResult, error)
}

func (m MangaServiceStub) Get(ctx context.Context, id string) (models.Manga, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, id)
	}
	return m.MangaService.Get(ctx, id)
}

func (m MangaServiceStub) Find(ctx context.Context, q models.MangaSearchQuery) (models.PaginatedMangaResult, error) {
	if m.FindFn != nil {
		return m.FindFn(ctx, q)
	}
	return m.MangaService.Find(ctx, q)
}

// ---- User service stub ----
type UserServiceStub struct {
	users.UserLibraryService

	UpdateFn func(ctx context.Context, userID, mangaID string, chapter int) error
}

func (u UserServiceStub) UpdateUserReadingProgress(ctx context.Context, userID, mangaID string, chapter int) error {
	if u.UpdateFn != nil {
		return u.UpdateFn(ctx, userID, mangaID, chapter)
	}
	return u.UserLibraryService.UpdateUserReadingProgress(ctx, userID, mangaID, chapter)
}

// ---- TCP Progress client stub ----
type TCPClientStub struct {
	client.ProgressSync
	Sent []update.ProgressUpdate
}

func (c *TCPClientStub) Send(u update.ProgressUpdate) {
	c.Sent = append(c.Sent, u)
}

// ────────────────────────────────────────────────
//
//	HELPER — START GRPC SERVER FOR TESTS
//
// ────────────────────────────────────────────────
func startGrpcTestServer(t *testing.T, srv mangapb.MangaServiceServer) (addr string, stop func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0") // random available port
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	mangapb.RegisterMangaServiceServer(grpcServer, srv)

	go grpcServer.Serve(lis)

	return lis.Addr().String(), grpcServer.Stop
}

//
// ────────────────────────────────────────────────
//   TESTS
// ────────────────────────────────────────────────
//

// ----------------------
// GetManga Success
// ----------------------
func TestGetManga_Success(t *testing.T) {
	mangaStub := MangaServiceStub{
		GetFn: func(ctx context.Context, id string) (models.Manga, error) {
			return models.Manga{ID: "m1", Title: "Naruto", TotalChapters: 700}, nil
		},
	}

	userStub := UserServiceStub{}
	tcpStub := &TCPClientStub{}

	server := mygrpc.NewMangaServiceGrpcServer(userStub, mangaStub, tcpStub)

	addr, stop := startGrpcTestServer(t, server)
	defer stop()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := mangapb.NewMangaServiceClient(conn)

	resp, err := client.GetManga(context.Background(), &mangapb.GetMangaRequest{Id: "m1"})
	require.NoError(t, err)
	require.Equal(t, "Naruto", resp.Title)
	require.Equal(t, int32(700), resp.TotalChapters)
}

// ----------------------
// GetManga Not Found
// ----------------------
func TestGetManga_NotFound(t *testing.T) {
	mangaStub := MangaServiceStub{
		GetFn: func(ctx context.Context, id string) (models.Manga, error) {
			return models.Manga{}, mangas.ErrMangaNotExistInDatabase
		},
	}
	userStub := UserServiceStub{}
	tcpStub := &TCPClientStub{}

	server := mygrpc.NewMangaServiceGrpcServer(userStub, mangaStub, tcpStub)

	addr, stop := startGrpcTestServer(t, server)
	defer stop()

	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := mangapb.NewMangaServiceClient(conn)

	_, err := client.GetManga(context.Background(), &mangapb.GetMangaRequest{Id: "unknown"})
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, st.Code())
}

// ----------------------
// SearchManga Success
// ----------------------
func TestSearchManga_Success(t *testing.T) {
	mangaStub := MangaServiceStub{
		FindFn: func(ctx context.Context, q models.MangaSearchQuery) (models.PaginatedMangaResult, error) {
			return models.PaginatedMangaResult{
				Total:  1,
				Limit:  10,
				Offset: 0,
				Results: []models.MangaListItem{
					{ID: "m1", Title: "One Piece", TotalChapters: 1100},
				},
			}, nil
		},
	}

	server := mygrpc.NewMangaServiceGrpcServer(UserServiceStub{}, mangaStub, &TCPClientStub{})

	addr, stop := startGrpcTestServer(t, server)
	defer stop()

	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := mangapb.NewMangaServiceClient(conn)

	resp, err := client.SearchManga(context.Background(), &mangapb.SearchRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "One Piece", resp.Results[0].Title)
}

// ----------------------
// UpdateProgress Success
// ----------------------
func TestUpdateProgress_Success(t *testing.T) {
	userStub := UserServiceStub{
		UpdateFn: func(ctx context.Context, userID, mangaID string, chapter int) error {
			return nil // success
		},
	}

	tcpStub := &TCPClientStub{}

	server := mygrpc.NewMangaServiceGrpcServer(userStub, MangaServiceStub{}, tcpStub)

	addr, stop := startGrpcTestServer(t, server)
	defer stop()

	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := mangapb.NewMangaServiceClient(conn)

	resp, err := client.UpdateProgress(context.Background(), &mangapb.ProgressRequest{
		UserId:  "u1",
		MangaId: "m1",
		Chapter: 10,
	})

	require.NoError(t, err)
	require.True(t, resp.Ok)

	// Validate that TCP stub received the update
	require.Len(t, tcpStub.Sent, 1)
	require.Equal(t, "u1", tcpStub.Sent[0].UserID)
	require.Equal(t, 10, tcpStub.Sent[0].Chapter)
}

// ----------------------
// UpdateProgress Invalid Chapter
// ----------------------
func TestUpdateProgress_InvalidChapter(t *testing.T) {
	userStub := UserServiceStub{
		UpdateFn: func(ctx context.Context, userID, mangaID string, chapter int) error {
			return mangas.ErrInValidCurrentChapter
		},
	}

	server := mygrpc.NewMangaServiceGrpcServer(userStub, MangaServiceStub{}, &TCPClientStub{})

	addr, stop := startGrpcTestServer(t, server)
	defer stop()

	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := mangapb.NewMangaServiceClient(conn)

	_, err := client.UpdateProgress(context.Background(), &mangapb.ProgressRequest{
		UserId:  "u1",
		MangaId: "m1",
		Chapter: -5,
	})

	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
