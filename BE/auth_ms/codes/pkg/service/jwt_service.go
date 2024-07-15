package service

import (
	"auth_ms/pkg/config"
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/helper/safeasync"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Claims struct {
	UserId   uint    `json:"user_id"`
	UserRole string  `json:"user_role"`
	TokenId  *string `json:"token_id"`
	Exp      *int64  `json:"exp"`
}

func SetClaims(userId uint, userRole string, tokenId *string) *Claims {
	return &Claims{
		UserId:   userId,
		UserRole: userRole,
		TokenId:  tokenId,
	}
}

// Encode data to base64 URL encoding
func base64URLEncode(data *[]byte) *string {
	encoded := strings.TrimRight(base64.URLEncoding.EncodeToString(*data), "=")
	return &encoded
}

// Decode base64 URL encoding data
func base64URLDecode(data *string) (*[]byte, error) {
	paddedData := *data + strings.Repeat("=", (4-len(*data)%4)%4)
	decoded, err := base64.URLEncoding.DecodeString(paddedData)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}

// Generate HMAC SHA256 signature
func hmacSHA256(data *string, secret *string) *string {
	h := hmac.New(sha256.New, []byte(*secret))
	h.Write([]byte(*data))
	hSum := h.Sum(nil)
	hSumP := &hSum
	signature := base64URLEncode(hSumP)
	return signature
}

// Generate JWT token
func generateJWT(claimsP *Claims) (*string, *int64, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return nil, nil, err
	}
	encodedHeader := base64URLEncode(&headerJSON)

	expirationTime, err := time.ParseDuration(config.GetJwtConfig().JwtExpiresIn)
	if err != nil {
		return nil, nil, err
	}

	expirationTimeUnix := time.Now().Add(expirationTime).Unix()
	claimsP.Exp = &expirationTimeUnix

	claimsJSON, err := json.Marshal(claimsP)
	if err != nil {
		return nil, nil, err
	}

	encodedClaims := base64URLEncode(&claimsJSON)

	unsignedToken := *encodedHeader + "." + *encodedClaims
	signature := hmacSHA256(&unsignedToken, &config.GetJwtConfig().JwtSecret)
	token := unsignedToken + "." + *signature

	return &token, &expirationTimeUnix, nil
}

// Generate Refresh Token
func generateRefreshToken(claimsP *Claims) (*string, *int64, error) {
	expirationTime, err := time.ParseDuration(config.GetJwtConfig().RefreshExpiresIn)
	if err != nil {
		return nil, nil, err
	}

	expirationTimeUnix := time.Now().Add(expirationTime).Unix()
	claimsP.Exp = &expirationTimeUnix

	claimsJSON, err := json.Marshal(claimsP)
	if err != nil {
		return nil, nil, err
	}

	encodedClaims := base64URLEncode(&claimsJSON)

	signature := hmacSHA256(encodedClaims, &config.GetJwtConfig().JwtSecret)
	refreshToken := *encodedClaims + "." + *signature

	return &refreshToken, &expirationTimeUnix, nil
}

func IssueJwtWithRefreshToken(userId uint, userRole string, tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	claimsP := SetClaims(userId, userRole, tokenIdP)
	results := make(chan *dto.TokenDto, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	safeasync.Run(func() {
		defer wg.Done() // Decrement the counter when the goroutine completes
		jwtTokenP, jwtExpiresInP, _ := generateJWT(claimsP)
		jwtTokenDataP := &dto.TokenDto{
			Type:     "JWT",
			Token:    jwtTokenP,
			TokenExp: jwtExpiresInP,
		}
		results <- jwtTokenDataP
	})

	safeasync.Run(func() {
		defer wg.Done() // Decrement the counter when the goroutine completes
		refreshTokenP, refreshExpiresInP, _ := generateRefreshToken(claimsP)
		refreshTokenDataP := &dto.TokenDto{
			Type:     "REFRESH",
			Token:    refreshTokenP,
			TokenExp: refreshExpiresInP,
		}
		results <- refreshTokenDataP
	})

	wg.Wait()
	close(results)

	tokenDataP := &dto.TokenDataDto{}
	for result := range results {
		if result.Type == "JWT" {
			tokenDataP.Jwt = result
		} else {
			tokenDataP.Refresh = result
		}
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: tokenDataP}, nil
}

// Validate JWT token
func VerifyJWT(token *string) (*response.GenericServiceResponseDto, error) {
	parts := strings.Split(*token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	encodedHeader := parts[0]
	encodedClaims := parts[1]
	providedSignature := parts[2]

	unsignedToken := encodedHeader + "." + encodedClaims
	expectedSignature := hmacSHA256(&unsignedToken, &config.GetJwtConfig().JwtSecret)

	if providedSignature != *expectedSignature {
		return nil, fmt.Errorf("invalid token signature")
	}

	claimsJSON, err := base64URLDecode(&encodedClaims)
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(*claimsJSON, &claims); err != nil {
		return nil, err
	}

	if time.Now().Unix() > *claims.Exp {
		return nil, fmt.Errorf("token has expired")
	}

	// return &claims, nil
	return &response.GenericServiceResponseDto{StatusCode: 200, Data: &claims}, nil
}

// Validate Refresh Token
func VerifyRefreshToken(token *string) (*response.GenericServiceResponseDto, error) {
	parts := strings.Split(*token, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	encodedClaims := parts[0]
	providedSignature := parts[1]

	expectedSignature := hmacSHA256(&encodedClaims, &config.GetJwtConfig().JwtSecret)

	if providedSignature != *expectedSignature {
		return nil, fmt.Errorf("invalid token signature")
	}

	claimsJSON, err := base64URLDecode(&encodedClaims)
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(*claimsJSON, &claims); err != nil {
		return nil, err
	}

	if time.Now().Unix() > *claims.Exp {
		return nil, fmt.Errorf("token has expired")
	}

	// return &claims, nil
	return &response.GenericServiceResponseDto{StatusCode: 200, Data: &claims}, nil
}
