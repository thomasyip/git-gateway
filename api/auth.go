package api

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/okta/okta-jwt-verifier-golang"
)

// requireAuthentication checks incoming requests for tokens presented using the Authorization header
func (a *API) requireAuthentication(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	logrus.Info("Getting auth token")
	token, err := a.extractBearerToken(w, r)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Parsing JWT claims: %v", token)
	return a.parseJWTClaims(token, r)
}

func (a *API) extractBearerToken(w http.ResponseWriter, r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", unauthorizedError("This endpoint requires a Bearer token")
	}

	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", unauthorizedError("This endpoint requires a Bearer token")
	}

	return matches[1], nil
}

func (a *API) parseJWTClaims(bearer string, r *http.Request) (context.Context, error) {
	// Reimplemented to use Okta lib
	// Original validation only work for HS256 algo,
	// Okta supports RS256 only which requires public key downloading and caching (key rotation)
	config := getConfig(r.Context())

	toValidate := map[string]string{}
	toValidate["aud"] = config.JWT.AUD
	toValidate["cid"] = config.JWT.CID

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer: config.JWT.Issuer,
		ClaimsToValidate: toValidate,
	}

	verifier := jwtVerifierSetup.New()

	_, err := verifier.VerifyAccessToken(bearer)

	// @TODO? WARNING: Should be roles and other claims be checked here?

	if err != nil {
		return nil, unauthorizedError("Invalid token: %v", err)
	}

	logrus.Infof("parseJWTClaims passed")

	// return nil, because the `github.go` is coded to send personal token
	// both github oauth generates its own id, so oauth pass-thru is impossible
	// we can improve the gateway to talk oauth with github.com, but we will
	// still return nil here.
	return nil, nil
}
