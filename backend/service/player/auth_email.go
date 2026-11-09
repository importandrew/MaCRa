package player

import (
	"context"
	"errors"
	"time"

	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

// SignUp creates the account and mails a verification link. For an email that
// already has an account it answers the same and mails that address a
// password reset link instead, so the response doesn't reveal which emails
// are registered.
func (s *Service) SignUp(ctx context.Context, email, name, password, origin string) error {
	player, err := s.Add(ctx, email, name, password)
	switch {
	case errors.Is(err, models.ErrPlayerExists):
		s.sendAccountExists(ctx, email, origin)
		return nil
	case err != nil:
		return err
	}

	s.sendVerification(ctx, player.ID, email, origin)
	return nil
}

func (s *Service) ResendVerification(ctx context.Context, email, origin string) error {
	player, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil //nolint:nilerr // don't reveal whether the email is registered
	}

	verified, err := s.repo.IsEmailVerified(ctx, player.ID)
	if err != nil || verified {
		return nil //nolint:nilerr // don't reveal the account state
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
		return nil //nolint:nilerr // don't reveal whether the email is registered
	}

	s.sendPasswordReset(ctx, player.ID, email, origin)
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, rawToken, password string) error {
	if err := checkNewPassword(password); err != nil {
		return err
	}

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
	s.sendAuthMail(ctx, userID, email, origin, models.AuthPurposeEmailVerify, models.EmailVerifyTTL, s.mail.SendVerification)
}

func (s *Service) sendPasswordReset(ctx context.Context, userID int, email, origin string) {
	s.sendAuthMail(ctx, userID, email, origin, models.AuthPurposePasswordReset, models.PasswordResetTTL, s.mail.SendPasswordReset)
}

// sendAccountExists tells the owner of an address that someone tried to sign
// up with it, with a password reset link.
func (s *Service) sendAccountExists(ctx context.Context, email, origin string) {
	player, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		logger.Errorf("account exists mail to %s: %v", email, err)
		return
	}

	s.sendAuthMail(ctx, player.ID, email, origin, models.AuthPurposePasswordReset, models.PasswordResetTTL, s.mail.SendAccountExists)
}

// sendAuthMail issues a token and mails it with send. Nothing is reported to
// the caller: the response is the same whether the mail went out or not. Over
// the address's limit nothing is sent.
func (s *Service) sendAuthMail(
	ctx context.Context,
	userID int,
	email, origin string,
	purpose models.AuthPurpose,
	ttl time.Duration,
	send func(to, origin, token string) error,
) {
	raw, err := s.issueToken(ctx, userID, purpose, ttl)
	if errors.Is(err, models.ErrAuthMailLimit) {
		logger.Warnf("%s mail to %s not sent: %v", purpose, email, err)
		return
	}
	if err != nil {
		logger.Errorf("issue %s token for %s: %v", purpose, email, err)
		return
	}

	if err := send(email, origin, raw); err != nil {
		logger.Errorf("send %s mail to %s: %v", purpose, email, err)
	}
}
