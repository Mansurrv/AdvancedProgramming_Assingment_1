package client

import (
	"context"
	"fmt"
	"time"

	"appointment-service/internal/apperr"
	doctorpb "doctor-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DoctorGRPCClient struct {
	client doctorpb.DoctorServiceClient
}

func NewDoctorGRPCClient(client doctorpb.DoctorServiceClient) *DoctorGRPCClient {
	return &DoctorGRPCClient{
		client: client,
	}
}

func (c *DoctorGRPCClient) DoctorExists(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := c.client.GetDoctor(ctx, &doctorpb.GetDoctorRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return false, err
		}

		switch st.Code() {
		case codes.NotFound:
			return false, nil
		case codes.Unavailable, codes.DeadlineExceeded:
			return false, apperr.ErrDoctorServiceUnavailable
		default:
			return false, fmt.Errorf("%w: %s", apperr.ErrDoctorServiceError, st.Message())
		}
	}

	return true, nil
}