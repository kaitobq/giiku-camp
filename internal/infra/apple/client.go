package apple

import (
	"context"
	"fmt"
	"os"

	"github.com/Timothylock/go-signin-with-apple/apple"
)

type Client interface {
	VerifyAuthorizationCode(ctx context.Context, code string) (*VerifyAuthorizationCodeOutput, error)
}

type client struct {
	cli          *apple.Client
	clientID     string
	clientSecret string
}

func New() (Client, error) {
	var (
		signingKey = os.Getenv("APPLE_SIGNING_KEY")
		teamID     = os.Getenv("APPLE_TEAM_ID")
		clientID   = os.Getenv("APPLE_CLIENT_ID")
		keyID      = os.Getenv("APPLE_KEY_ID")
	)

	clientSecret, err := apple.GenerateClientSecret(signingKey, teamID, clientID, keyID)
	if err != nil {
		return nil, err
	}
	cli := apple.New()

	return &client{
		cli:          cli,
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

type VerifyAuthorizationCodeOutput struct {
	RefreshToken string
	UserID       string
}

func (c *client) VerifyAuthorizationCode(ctx context.Context, code string) (*VerifyAuthorizationCodeOutput, error) {
	var resp apple.ValidationResponse
	if err := c.cli.VerifyAppToken(ctx, apple.AppValidationTokenRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Code:         code,
	}, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("VerifyAppToken  %s - %s", resp.Error, resp.ErrorDescription)
	}
	userID, err := apple.GetUniqueID(resp.IDToken)
	if err != nil {
		return nil, err
	}

	return &VerifyAuthorizationCodeOutput{
		RefreshToken: resp.RefreshToken,
		UserID:       userID,
	}, nil
}
