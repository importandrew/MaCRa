package player

import (
	"context"

	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

func (s *Service) SignUp(ctx context.Context, email, name, password, origin string) error {
	player, err := s.Add(ctx, email, name, password)
	if err != nil {
		return err
	}

	s.sendVerification(ctx, player.ID, email, origin)
	return nil
}

func (s *Service) ResendVerification(ctx context.Context, email, origin string) error {
	player, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil
	}

	verified, err := s.repo.IsEmailVerified(ctx, player.ID)
	if err != nil || verified {
		return nil
	}

	s.sendVerification(ctx, player.ID, email, origin)
	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, rawToken string) error {
	token, err := s.lookupToken(ctx, rawToken, models.AuthPurposeEmailVerify)
	if err != nil {
		return err
	}

	if err := s.repo.SetEmailVerified(ctx, token.UserID); err != nil {
		return err
	}

	if err := s.repo.MarkAuthTokenUsed(ctx, token.ID); err != nil {
		return err
	}

	return nil
}

func (s *Service) ForgotPassword(ctx context.Context, email, origin string) error {
	player, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil
	}

	s.sendPasswordReset(ctx, player.ID, email, origin)
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, rawToken, password string) error {
	token, err := s.lookupToken(ctx, rawToken, models.AuthPurposePasswordReset)
	if err != nil {
		return err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, token.UserID, hash); err != nil {
		return err
	}

	if err := s.repo.MarkAuthTokenUsed(ctx, token.ID); err != nil {
		return err
	}

	_ = s.repo.InvalidateAuthTokens(ctx, token.UserID, models.AuthPurposePasswordReset)
	return nil
}

func (s *Service) sendVerification(ctx context.Context, userID int, email, origin string) {
	raw, err := s.issueToken(ctx, userID, models.AuthPurposeEmailVerify, models.EmailVerifyTTL)
	if err != nil {
		logger.Errorf("issue verify token for %s: %v", email, err)
		return
	}

	if err := s.mail.SendVerification(email, origin, raw); err != nil {
		logger.Errorf("send verify email to %s: %v", email, err)
	}
}

func (s *Service) sendPasswordReset(ctx context.Context, userID int, email, origin string) {
	raw, err := s.issueToken(ctx, userID, models.AuthPurposePasswordReset, models.PasswordResetTTL)
	if err != nil {
		logger.Errorf("issue reset token for %s: %v", email, err)
		return
	}

	if err := s.mail.SendPasswordReset(email, origin, raw); err != nil {
		logger.Errorf("send reset email to %s: %v", email, err)
	}
}
