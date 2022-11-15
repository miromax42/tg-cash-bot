package gapi

import (
	"context"
	"strconv"

	"github.com/cockroachdb/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/pb"
)

type chatID int64

func (r chatID) Recipient() string {
	return strconv.FormatInt(int64(r), 10)
}

func (s *Server) SendReport(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	_, err := s.t.Send(chatID(req.UserId), req.Message)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			errors.Wrapf(err, "cant send to user=%d", req.UserId).Error())
	}

	return &pb.SendMessageResponse{Success: true}, nil
}
