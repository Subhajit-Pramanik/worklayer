package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
	"google.golang.org/grpc/metadata"
)

func IamJWTVerify(authClient iAMV1.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractIamJWTFormCookie(c)
		if token == "" {
			token = extractIamJWTFormHeader(c)
		}

		if token == "" {
			return errors.Unauthorized("Auth token is required")
		}

		resp, err := authClient.ValidateSession(c.UserContext(), &iAMV1.ValidateSessionRequest{
			AccessToken: token,
		})
		if err != nil || !resp.GetIsValid() {
			return errors.Unauthorized("invalid, expired, or revoked auth token")
		}

		user := resp.GetUser()
		if user == nil {
			return errors.Unauthorized("invalid or expired auth token")
		}

		ctx := ctxutil.InjectIAMUserID(c.UserContext(), user.GetId())
		ctx = ctxutil.InjectIAMUserEmail(ctx, user.GetEmail())

		ctx = metadata.AppendToOutgoingContext(ctx,
			"iam_user_id", user.GetId(),
			"iam_user_email", user.GetEmail(),
		)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func extractIamJWTFormCookie(c *fiber.Ctx) string {
	return c.Cookies("__vyo_iam_auth")
}

func extractIamJWTFormHeader(c *fiber.Ctx) string {
	str := c.Get("Authorization")
	if str == "" {
		return ""
	}

	parts := strings.Split(str, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
