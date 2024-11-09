package services

import (
	"AuthAPI/cfg"
	"AuthAPI/internal/services/auth"
	"AuthAPI/internal/util"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthService struct {
	oauthCfg  *oauth2.Config
	jwtSecret string
}

func NewGoogleAuthService(cfg *cfg.AppConfig) *GoogleAuthService {
	return &GoogleAuthService{
		oauthCfg: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		jwtSecret: cfg.JwtSecret,
	}
}

func (service GoogleAuthService) GetAuthUrl(c *fiber.Ctx) error {
	authUrl := service.oauthCfg.AuthCodeURL("state", oauth2.AccessTypeOffline)

	return c.JSON(fiber.Map{"authUrl": authUrl})
}

func (service GoogleAuthService) HandleCallBack(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Code not provided"})
	}

	token, err := service.oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to exchange code for token"})
	}

	tokenSource := service.oauthCfg.TokenSource(context.Background(), token)
	oauth2Service := oauth2.NewClient(context.Background(), tokenSource)

	response, err := oauth2Service.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to fetch userInfo from google"})
	}

	defer response.Body.Close()

	userInfo := auth.GoogleUserInfo{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user info"})
	}

	jwt, err := util.GenerateTokenFromGoogleIdentity(userInfo, service.jwtSecret)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JWT token"})
	}

	return c.JSON(fiber.Map{"token": jwt})

}
