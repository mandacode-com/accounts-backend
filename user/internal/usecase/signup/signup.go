package signup

import (
	"context"

	"github.com/google/uuid"
	authrepo "mandacode.com/accounts/user/internal/repository/auth"
	authrepodto "mandacode.com/accounts/user/internal/repository/auth/dto"
	dbrepo "mandacode.com/accounts/user/internal/repository/database"
	profilerepo "mandacode.com/accounts/user/internal/repository/profile"
	profilerepodto "mandacode.com/accounts/user/internal/repository/profile/dto"
	usereventrepo "mandacode.com/accounts/user/internal/repository/userevent"
	signupdto "mandacode.com/accounts/user/internal/usecase/signup/dto"
)

type SingupUsecase struct {
	authRepo              *authrepo.AuthRepository
	profileRepo           *profilerepo.ProfileRepository
	dbUserRepo            *dbrepo.UserRepository
	userEventEmitter      *usereventrepo.UserEventEmitter
	emailVerificationLink string
}

// NewSignupUsecase creates a new instance of LocalSingupUsecase with the provided repositories.
func NewSignupUsecase(
	authRepo *authrepo.AuthRepository,
	profileRepo *profilerepo.ProfileRepository,
	dbUserRepo *dbrepo.UserRepository,
	userEventEmitter *usereventrepo.UserEventEmitter,
) *SingupUsecase {
	return &SingupUsecase{
		authRepo:         authRepo,
		profileRepo:      profileRepo,
		dbUserRepo:       dbUserRepo,
		userEventEmitter: userEventEmitter,
	}
}

// LocalSignup performs the signup process for a new user.
func (s *SingupUsecase) LocalSignup(ctx context.Context, req *signupdto.LocalSignupRequest) (*signupdto.LocalSignupResponse, error) {
	userID := uuid.New()
	dbUser, err := s.dbUserRepo.CreateUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.authRepo.CreateLocalUser(ctx, &authrepodto.CreateLocalUserRequest{
		UserID:   userID,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if dbErr := s.dbUserRepo.DeleteUser(ctx, userID); dbErr != nil {
			return nil, dbErr
		}
		if err := s.userEventEmitter.EmitUserDeletedEvent(ctx, userID); err != nil {
			return nil, err
		}
		return nil, err
	}
	_, err = s.profileRepo.CreateProfileUser(ctx, &profilerepodto.CreateProfileUserRequest{
		UserID:   userID,
		Email:    req.Email,
		SyncCode: dbUser.SyncCode,
	})
	if err != nil {
		if dbErr := s.dbUserRepo.DeleteUser(ctx, userID); dbErr != nil {
			return nil, dbErr
		}
		if err := s.userEventEmitter.EmitUserDeletedEvent(ctx, userID); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &signupdto.LocalSignupResponse{
		UserID:    userID,
		Email:     req.Email,
		CreatedAt: dbUser.CreatedAt,
	}, nil
}

// OAuthSignup performs the OAuth signup process for a new user.
func (s *SingupUsecase) OAuthSignup(ctx context.Context, req *signupdto.OAuthSignupRequest) (*signupdto.OAuthSignupResponse, error) {
	userID := uuid.New()
	dbUser, err := s.dbUserRepo.CreateUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	authUser, err := s.authRepo.CreateOAuthUser(ctx, &authrepodto.CreateOAuthUserRequest{
		UserID:      userID,
		Provider:    req.Provider,
		AccessToken: &req.AccessToken,
	})
	if err != nil {
		if dbErr := s.dbUserRepo.DeleteUser(ctx, userID); dbErr != nil {
			return nil, dbErr
		}
		if err := s.userEventEmitter.EmitUserDeletedEvent(ctx, userID); err != nil {
			return nil, err
		}
		return nil, err
	}
	_, err = s.profileRepo.CreateProfileUser(ctx, &profilerepodto.CreateProfileUserRequest{
		UserID:   userID,
		Email:    authUser.Email,
		SyncCode: dbUser.SyncCode,
	})
	if err != nil {
		if dbErr := s.dbUserRepo.DeleteUser(ctx, userID); dbErr != nil {
			return nil, dbErr
		}
		if err := s.userEventEmitter.EmitUserDeletedEvent(ctx, userID); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &signupdto.OAuthSignupResponse{
		UserID:    userID,
		Provider:  req.Provider,
		Email:     authUser.Email,
		CreatedAt: dbUser.CreatedAt,
	}, nil
}
