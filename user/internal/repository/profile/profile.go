package profilerepo

import (
	"context"

	profilev1 "github.com/mandacode-com/accounts-proto/go/profile/v1"
	profilerepodto "mandacode.com/accounts/user/internal/repository/profile/dto"
)

type ProfileRepository struct {
	client profilev1.ProfileServiceClient
}

// NewProfileRepository creates a new instance of ProfileRepository with the provided client.
func NewProfileRepository(client profilev1.ProfileServiceClient) *ProfileRepository {
	return &ProfileRepository{
		client: client,
	}
}

func (p *ProfileRepository) CreateProfileUser(ctx context.Context, req *profilerepodto.CreateProfileUserRequest) (*profilerepodto.CreateProfileUserResponse, error) {
	protoRes, err := p.client.InitUser(ctx, req.ToProto())
	if err != nil {
		return nil, err
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, err
	}
	res, err := profilerepodto.NewCreateProfileUserResponse(protoRes)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ProfileRepository) UpdateEmail(ctx context.Context, req *profilerepodto.UpdateEmailRequest) (*profilerepodto.UpdateEmailResponse, error) {
	protoRes, err := p.client.UpdateEmail(ctx, req.ToProto())
	if err != nil {
		return nil, err
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, err
	}
	res, err := profilerepodto.NewUpdateEmailResponse(protoRes)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ProfileRepository) DeleteProfileUser(ctx context.Context, req *profilerepodto.DeleteProfileUserRequest) (*profilerepodto.DeleteProfileUserResponse, error) {
	protoRes, err := p.client.DeleteUser(ctx, req.ToProto())
	if err != nil {
		return nil, err
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, err
	}
	res, err := profilerepodto.NewDeleteProfileUserResponse(protoRes)
	if err != nil {
		return nil, err
	}
	return res, nil
}
