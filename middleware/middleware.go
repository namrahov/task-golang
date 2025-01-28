package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
	"task-golang/config"
	"task-golang/model"
	"task-golang/repo"
)

var headers = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
	"DP-Customer-ID",
	"DP-User-ID",
	"User-Agent",
	"X-Forwarded-For",
	"requestid",
}

// Define the whitelist
var whitelist = map[string][]string{
	"GET":  {"/v1/users/active", "/swagger/*"},
	"POST": {"/v1/users/login", "/v1/users/register"},
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("siledi0")
			ctx := r.Context()

			requestID := r.Header.Get(model.HeaderKeyRequestID)
			operation := r.RequestURI
			customerID := r.Header.Get(model.HeaderKeyCustomerID)
			userID := r.Header.Get(model.HeaderKeyUserID)
			userAgent := r.Header.Get(model.HeaderKeyUserAgent)
			userIP := r.Header.Get(model.HeaderKeyUserIP)

			if len(requestID) == 0 {
				requestID = uuid.New().String()
			}
			fields := log.Fields{}
			addLoggerParam(fields, model.LoggerKeyRequestID, requestID)
			addLoggerParam(fields, model.LoggerKeyCustomerID, customerID)
			addLoggerParam(fields, model.LoggerKeyOperation, operation)
			addLoggerParam(fields, model.LoggerKeyUserAgent, userAgent)
			addLoggerParam(fields, model.LoggerKeyUserID, userID)
			addLoggerParam(fields, model.LoggerKeyUserIP, userIP)

			logger := log.WithFields(fields)
			header := http.Header{}

			for _, v := range headers {
				header.Add(v, r.Header.Get(v))
			}

			// Check if the request URL and method are in the whitelist
			fmt.Println("isledi1", r.Method, " url = ", r.URL.Path)
			if isWhitelisted(r.Method, r.URL.Path) {
				// Proceed to the next handler
				ctx = context.WithValue(ctx, model.ContextLogger, logger)
				ctx = context.WithValue(ctx, model.ContextHeader, header)
				fmt.Println("isledi2")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Extract Authorization Header
			authHeader := r.Header.Get("Authorization")
			fmt.Println("authHeader", authHeader)
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") || !existByToken(ctx, authHeader) {
				http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
				return
			}

			// Extract Token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and Validate JWT Token
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(config.Props.JwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			// Extract Claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims["user_id"] == nil || claims["roles"] == nil {
				http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
				return
			}

			userId := int64(claims["user_id"].(float64))
			roles := claims["roles"].([]interface{})

			userRoles := make([]string, 0, len(roles)) // Initialize with capacity but flexible length

			for i, role := range roles {
				fmt.Printf("Processing role at index %d: %v (type: %T)\n", i, role, role)
				// Ensure role is a map and extract the "name" field
				if roleMap, ok := role.(map[string]interface{}); ok {
					if name, exists := roleMap["name"]; exists {
						if nameStr, ok := name.(string); ok {
							userRoles = append(userRoles, nameStr)
						} else {
							fmt.Printf("name field at index %d is not a string: %v\n", i, name)
						}
					} else {
						fmt.Printf("No 'name' field found in role at index %d: %v\n", i, role)
					}
				} else {
					fmt.Printf("Role at index %d is not a map: %v\n", i, role)
				}
			}

			fmt.Println("isledi10", userRoles)
			hasPermission := checkPermission(userRoles, r.RequestURI, r.Method)
			if !hasPermission {
				http.Error(w, "Forbidden: You do not have access to this resource", http.StatusForbidden)
				return
			}

			fmt.Println("isledi11")
			// Add User Info to Context
			ctx = context.WithValue(ctx, model.ContextUserID, userId)
			ctx = context.WithValue(ctx, model.ContextUserRoles, userRoles)
			ctx = context.WithValue(ctx, model.ContextAuthHeader, authHeader)

			ctx = context.WithValue(ctx, model.ContextLogger, logger)
			ctx = context.WithValue(ctx, model.ContextHeader, header)

			fmt.Println("isledi12")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func addLoggerParam(fields log.Fields, field string, value string) {
	if len(value) > 0 {
		fields[field] = value
	}
}

func checkPermission(roles []string, requestURI, httpMethod string) bool {
	// Ensure roles is not empty to avoid SQL syntax errors
	if len(roles) == 0 {
		log.Errorf("Roles slice is empty, cannot check permissions.")
		return false
	}

	var permissions []model.Permission

	// Query permissions with GORM
	err := repo.Db.Table("permissions").
		Select("permissions.*").
		Joins("JOIN roles_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN roles r ON r.id = rp.role_id").
		Where("r.name IN ?", roles).
		Find(&permissions).Error
	if err != nil {
		fmt.Printf("checkPermission err: %v\n", err)
		log.Errorf("Error fetching permissions: %v", err)
		return false
	}

	// Check if the request URI and method match any of the permissions
	for _, permission := range permissions {
		if matchPattern(permission.URL, requestURI) && strings.EqualFold(permission.HTTPMethod, httpMethod) {
			return true
		}
	}
	return false
}

// Helper function to check if the method and URL are in the whitelist
func isWhitelisted(method, url string) bool {
	allowedURLs, exists := whitelist[method]
	if !exists {
		return false
	}
	for _, allowedURL := range allowedURLs {
		// Check if the allowed URL ends with a wildcard '*' or matches exactly
		if strings.HasSuffix(allowedURL, "*") {
			// Remove the '*' and check if the URL starts with the remaining prefix
			prefix := strings.TrimSuffix(allowedURL, "*")
			if strings.HasPrefix(url, prefix) {
				return true
			}
		} else if allowedURL == url {
			return true
		}
	}
	return false
}

// matchPattern checks if a request URI matches a permission pattern
func matchPattern(pattern, requestURI string) bool {
	// Remove query parameters (everything after '?')
	cleanedRequestURI := strings.Split(requestURI, "?")[0]

	// Replace placeholders in the pattern with regex patterns
	regexPattern := "^" + strings.ReplaceAll(pattern, "{id}", "\\d+") + "$"

	// Match the cleaned requestURI against the compiled regex
	matched, err := regexp.MatchString(regexPattern, cleanedRequestURI)
	if err != nil {
		log.Errorf("Error matching pattern: %v", err)
		return false
	}

	return matched
}

// ExistByToken checks if a token exists in Redis
func existByToken(ctx context.Context, token string) bool {
	// Construct the token index key
	tokenIndexKey := fmt.Sprintf("tokenIndex:%s", token)

	// Check if the key exists in Redis
	exists, err := repo.RedisClient.Exists(ctx, tokenIndexKey).Result()
	if err != nil {
		_ = fmt.Errorf("error checking if token exists in Redis: %w", err)
		return false
	}

	// Redis EXISTS command returns the number of keys that exist (0 or 1 in this case)
	fmt.Println(exists > 0)
	return exists > 0
}
