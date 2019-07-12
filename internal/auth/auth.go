package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	authenticationapi "k8s.io/api/authentication/v1"
	clientgo "k8s.io/client-go/kubernetes"
	"log"
	"time"
)

const audience = "aether"

type AetherUserClaims struct {
	ServiceAccountName string `json:"foo"`
	jwt.StandardClaims
}

var privKey = generatePrivateKey()

// TODO don't do this
func generatePrivateKey() []byte {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("failed to generate key")
		return nil
	}

	// Convert it to pem
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	fmt.Println("Generated private key")
	return pem.EncodeToMemory(block)
}

func validateKubernetesToken(tokenReviewResponse *authenticationapi.TokenReview) string {
	// Check if Kubernetes OK'd the token.
	if tokenReviewResponse.Status.Authenticated == true {

		// Check that we're the audience for the token.
		for _, v := range tokenReviewResponse.Status.Audiences {
			if v == audience {
				log.Println("Token match")
				return tokenReviewResponse.Status.User.Username
			}
		}

		log.Println("Not in audiences")
		return ""

	} else {
		log.Printf("authenticated field was false in bearerToken")
	}

	return ""
}

// Return the username of the caller.
// If the caller can't be authenticated, return the empty string.
func authenticateKubernetesToken(clientset *clientgo.Clientset, serviceToken string) string {
	res, err := clientset.AuthenticationV1().TokenReviews().Create(&authenticationapi.TokenReview{
		Spec: authenticationapi.TokenReviewSpec{
			Token: serviceToken,
		},
	})
	if err != nil {
		log.Println("Error checking Kubernetes token: ", err)
		return ""
	}

	return validateKubernetesToken(res)
}

func generateToken(user string) (string, error) {
	claims := AetherUserClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			NotBefore: time.Now().Add(-time.Second).Unix(), // 1 second leeway.
			Issuer:    "aether",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(privKey)
}

func Authenticate(tokenStr string) string {
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if err != nil {
		return ""
	}

	if parsedToken.Valid {
		return "username from token"
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Println("That's not even a parsedToken")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			log.Println("Timing is everything")
		} else {
			log.Println("Couldn't handle this parsedToken:", err)
		}
	} else {
		log.Println("Couldn't handle this parsedToken:", err)
	}

	return ""
}

// Authenticates using a Kubernetes service account token.
// If successful, returns a JWT from Aether. The Aether JWT should be used for all other authentication calls.
func InitialAuthentication(clientset *clientgo.Clientset, kubernetesServiceToken string) (string, error) {
	authFailed := errors.New("authentication failed")

	user := authenticateKubernetesToken(clientset, kubernetesServiceToken)
	if len(user) == 0 {
		return "", authFailed
	}

	aetherToken, err := generateToken(user)
	if err != nil {
		log.Println("Failed to generate token")
		return "", authFailed
	}

	return aetherToken, nil
}
