package service

import (
	"context"
	"fmt"
	"peer-link-server/config"
	"peer-link-server/internal/model"
	"peer-link-server/internal/repository"
	apperrors "peer-link-server/pkg/errors"
	"time"

	"github.com/livekit/protocol/auth"
	livekit "github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type LiveKitService interface {
	CreateRoom(ctx context.Context, req CreateRoomRequest) (*model.Room, error)
	GetToken(ctx context.Context, roomName, identity string) (*TokenResponse, error)
	ListRooms(ctx context.Context, page, pageSize int) ([]model.Room, int64, error)
	DeleteRoom(ctx context.Context, roomName string) error
	ListParticipants(ctx context.Context, roomName string) ([]*livekit.ParticipantInfo, error)
	RemoveParticipant(ctx context.Context, roomName, identity string) error
}

type CreateRoomRequest struct {
	Name            string `json:"name"             binding:"required,min=3,max=64"`
	DisplayName     string `json:"display_name"`
	MaxParticipants int    `json:"max_participants"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ServerURL string `json:"server_url"`
	RoomName  string `json:"room_name"`
	Identity  string `json:"identity"`
	ExpiresAt int64  `json:"expires_at"`
}

type liveKitService struct {
	cfg        *config.LiveKitConfig
	roomRepo   repository.RoomRepository
	roomClient *lksdk.RoomServiceClient
}

func NewLiveKitService(cfg *config.LiveKitConfig, roomRepo repository.RoomRepository) LiveKitService {
	client := lksdk.NewRoomServiceClient(cfg.Host, cfg.APIKey, cfg.APISecret)
	return &liveKitService{cfg: cfg, roomRepo: roomRepo, roomClient: client}
}

func (s *liveKitService) CreateRoom(ctx context.Context, req CreateRoomRequest) (*model.Room, error) {
	maxP := req.MaxParticipants
	if maxP <= 0 {
		maxP = 20
	}
	// 在 LiveKit Server 创建房间（设置参数）
	_, err := s.roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:            req.Name,
		EmptyTimeout:    300,
		MaxParticipants: uint32(maxP),
	})
	if err != nil {
		return nil, apperrors.Wrap(500, 50001, "livekit create room failed", err)
	}
	// 本地数据库记录
	room := &model.Room{
		Name:            req.Name,
		DisplayName:     req.DisplayName,
		MaxParticipants: maxP,
	}
	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *liveKitService) GetToken(ctx context.Context, roomName, identity string) (*TokenResponse, error) {
	// 验证房间存在
	if _, err := s.roomRepo.FindByName(ctx, roomName); err != nil {
		return nil, err
	}
	ttl := s.cfg.TokenTTL
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	at := auth.NewAccessToken(s.cfg.APIKey, s.cfg.APISecret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   boolPtr(true),
		CanSubscribe: boolPtr(true),
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(ttl)

	token, err := at.ToJWT()
	if err != nil {
		return nil, apperrors.Wrap(500, 50002, "generate token failed", err)
	}
	return &TokenResponse{
		Token:     token,
		ServerURL: fmt.Sprintf("wss://%s", s.cfg.Host),
		RoomName:  roomName,
		Identity:  identity,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}, nil
}

func (s *liveKitService) ListRooms(ctx context.Context, page, pageSize int) ([]model.Room, int64, error) {
	return s.roomRepo.List(ctx, page, pageSize)
}

func (s *liveKitService) DeleteRoom(ctx context.Context, roomName string) error {
	_, err := s.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: roomName})
	if err != nil {
		return apperrors.Wrap(500, 50003, "livekit delete room failed", err)
	}
	return s.roomRepo.Delete(ctx, roomName)
}

func (s *liveKitService) ListParticipants(ctx context.Context, roomName string) ([]*livekit.ParticipantInfo, error) {
	res, err := s.roomClient.ListParticipants(ctx, &livekit.ListParticipantsRequest{Room: roomName})
	if err != nil {
		return nil, apperrors.Wrap(500, 50004, "livekit list participants failed", err)
	}
	return res.Participants, nil
}

func (s *liveKitService) RemoveParticipant(ctx context.Context, roomName, identity string) error {
	_, err := s.roomClient.RemoveParticipant(ctx, &livekit.RoomParticipantIdentity{
		Room: roomName, Identity: identity,
	})
	return err
}

func boolPtr(b bool) *bool { return &b }
