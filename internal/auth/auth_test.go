package auth

import (
	authv1 "k8s.io/api/authentication/v1"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		validation *authv1.TokenReview
		expect     string
	}{
		{
			validation: &authv1.TokenReview{
				Status: authv1.TokenReviewStatus{
					Authenticated: true,
					Audiences:     []string{"api", "aether"},
					User: authv1.UserInfo{
						Username: "vallery",
					},
				},
			},
			expect: "vallery",
		},
		{
			validation: &authv1.TokenReview{},
			expect:     "",
		},
		{
			validation: &authv1.TokenReview{
				Status: authv1.TokenReviewStatus{
					Audiences: []string{"api", "aether"},
					User: authv1.UserInfo{
						Username: "vallery",
					},
				},
			},
			expect: "",
		},
		{
			validation: &authv1.TokenReview{
				Status: authv1.TokenReviewStatus{
					Authenticated: true,
					Audiences:     []string{"api"},
					User: authv1.UserInfo{
						Username: "vallery",
					},
				},
			},
			expect: "",
		},
		{
			validation: &authv1.TokenReview{
				Status: authv1.TokenReviewStatus{
					Authenticated: true,
					Audiences:     []string{"api", "aether"},
				},
			},
			expect: "",
		},
	}

	for _, c := range cases {
		actual := validateKubernetesToken(c.validation)
		if c.expect != actual {
			t.Errorf("Expected '%s', got '%s'", c.expect, actual)
		}
	}
}

func TestGenerateToken(t *testing.T) {
	cases := []struct {
		username string
	}{
		{
			username: "app",
		},
	}

	for _, c := range cases {
		_, err := generateToken(c.username)
		if err != nil {
			t.Errorf("Failed to generate token: %v", err)
		}
	}
}

func TestGeneratePrivateKey(t *testing.T) {
	key := generatePrivateKey()
	if len(key) == 0 {
		t.Errorf("key generation failed")
	}
}
