package goapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"cloud.google.com/go/spanner"
	firebase "firebase.google.com/go"
)

func (config *Config) Authenticate(req *http.Request, accountIDs []string) (authorizedAccountIDs []string, err error) {
	urlQuery := req.URL.Query()

	// Get API Key from headers or from URL parameters
	apiKey := req.Header.Get("X-API-KEY")
	if apiKey == "" {
		apiKey = urlQuery.Get("api_key")
	}

	if apiKey != "" {
		authorizedAccountIDs, err = config.AuthorizeApiKey(apiKey, accountIDs)
		if err != nil {
			return
		}
	} else {
		var userEmail string

		// Verify GIP ID Token and get the current user ID
		userEmail, err = VerifyToken(req.Header.Get("Authorization"))
		if err != nil {
			return
		}

		// Find account IDs and permissions for current user
		authorizedAccountIDs, err = config.AuthorizeToken(userEmail, accountIDs)
	}

	if len(authorizedAccountIDs) == 0 {
		err = errors.New("accounts not found")
	}

	return
}

// AuthorizeApiKey ...
func (config *Config) AuthorizeApiKey(apiKey string, accountIDs []string) (authorizedAccountIDs []string, err error) {
	query := `
		SELECT
			ApiKeys.AccountId as account_id
		FROM ApiKeys
		LEFT OUTER JOIN Roles ON Roles.RoleName = ApiKeys.RoleName
		WHERE ApiKeys.ApiKey = @api_key
			AND ApiKeys.AccountId IN UNNEST (@account_ids)
			AND LOWER(Roles.ServicePath) = @service_path
			AND LOWER(Roles.ServiceMethod) = @service_method
	`

	// Find all accounts if no account ids mentioned
	if len(accountIDs) == 0 {
		query = strings.Replace(query, "AND ApiKeys.AccountId IN UNNEST (@account_ids)", "", -1)
	}

	stmt := spanner.NewStatement(query)
	stmt.Params["api_key"] = apiKey
	stmt.Params["account_ids"] = accountIDs
	stmt.Params["service_path"] = strings.ToLower(config.Service.Path)
	stmt.Params["service_method"] = strings.ToLower(config.Service.Method)

	ctx := context.Background()
	iter := DB.Single().Query(ctx, stmt)
	defer iter.Stop()

	err = iter.Do(func(r *spanner.Row) error {
		var accountID string

		if err := r.Column(0, &accountID); err != nil {
			return err
		}
		authorizedAccountIDs = append(authorizedAccountIDs, accountID)

		return nil
	})

	return
}

func (config *Config) AuthorizeToken(userEmail string, accountIDs []string) (authorizedAccountIDs []string, err error) {
	query := `
		SELECT
			Users.AccountId as account_id
		FROM Users
		LEFT OUTER JOIN Roles ON Roles.RoleName = Users.RoleName
		WHERE Users.UserEmail = @user_email
			AND Users.AccountId IN UNNEST (@account_ids)
			AND LOWER(Roles.ServicePath) = @service_path
			AND LOWER(Roles.ServiceMethod) = @service_method
	`

	// Find all accounts if no account ids mentioned
	if len(accountIDs) == 0 {
		query = strings.Replace(query, "AND Users.AccountId IN UNNEST (@account_ids)", "", -1)
	}

	stmt := spanner.NewStatement(query)
	stmt.Params["user_email"] = userEmail
	stmt.Params["account_ids"] = accountIDs
	stmt.Params["service_path"] = strings.ToLower(config.Service.Path)
	stmt.Params["service_method"] = strings.ToLower(config.Service.Method)

	ctx := context.Background()
	iter := DB.Single().Query(ctx, stmt)
	defer iter.Stop()

	err = iter.Do(func(r *spanner.Row) error {
		var accountID string

		if err := r.Column(0, &accountID); err != nil {
			return err
		}
		authorizedAccountIDs = append(authorizedAccountIDs, accountID)

		return nil
	})

	return
}

// VerifyToken ...
func VerifyToken(jwtToken string) (userID string, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "  Execption in VerifyToken: ", string(debug.Stack()))
		}
	}()

	tokenParts := strings.Split(jwtToken, " ")
	if len(tokenParts) > 1 {
		jwtToken = strings.TrimSpace(tokenParts[1])
	} else {
		jwtToken = strings.TrimSpace(tokenParts[0])
	}

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Println("VerifyToken firebase.NewApp err: ", err)
		return "", err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Println("VerifyToken app.Auth err: ", err)
		return "", err
	}

	token, err := client.VerifyIDToken(ctx, jwtToken)
	if err != nil {
		log.Println("VerifyToken client.VerifyIDToken err: ", err)
		return "", err
	}

	// token.Claims : map[
	// 	auth_time:1.610686521e+09
	// 	email:someone@gmail.com
	// 	email_verified:false
	// 	firebase:map[
	// 		identities:map[
	// 			email:[someone@gmail.com]
	// 			phone:[+919293949596]
	// 		]
	// 		sign_in_provider:password
	// 	]
	// 	name:giri giri
	// 	phone_number:+919293949596
	// 	user_id:8SyEz1pSizc1IrIK36zigfy6Ou12
	// ]

	return fmt.Sprintf("%v", token.Claims["email"]), nil
}
