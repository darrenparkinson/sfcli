package salesforce

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type sftoken struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

// Client is the main salesforce client for interacting with the library.  It can be created using NewClient
type Client struct {
	// BaseURL for the API.  Set using `salesforce.New()`.
	BaseURL string

	// Version of the API to use.  Default is v53.0. Other versions haven't been tested.
	// Set the Version field after initialising with New.
	Version string

	//HTTP Client to use for making requests, allowing the user to supply their own if required.
	HTTPClient *http.Client

	// BulkService represents the Bulk 2.0 API
	BulkService *BulkService
	// AccountService represents the Account object
	AccountService *AccountService
	// ContactService represents the Contact object
	ContactService *ContactService
	// OpportunityService represents the Opportunity object
	OpportunityService *OpportunityService
	// UserService represents the User object
	UserService *UserService

	username string
	password string
	clientID string
	secret   string
	token    *sftoken
	lim      *rate.Limiter
}

// BulkService represents the Bulk Service 2.0 API
type BulkService struct {
	client *Client
}

// AccountService represents the Account object
type AccountService struct {
	client *Client
}

// ContactService represents the Contact object
type ContactService struct {
	client *Client
}

// OpportunityService represents the Opportunity object
type OpportunityService struct {
	client *Client
}

// UserService represents the User object
type UserService struct {
	client *Client
}

// NewClient is a helper function that returns an new salesforce client given the required parameters.
// Optionally you can provide your own http client or use nil to use the default.  This is done to
// ensure you're aware of the decision you're making to not provide your own http client.
func NewClient(baseURL, username, password, clientID, secret string, client *http.Client) (*Client, error) {
	if baseURL == "" || username == "" || password == "" || clientID == "" || secret == "" {
		return nil, errors.New("missing required parameters")
	}
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	rl := rate.NewLimiter(150, 1) // TODO: Identify what this should be
	c := &Client{
		BaseURL:    baseURL,
		Version:    "v53.0",
		HTTPClient: client,
		username:   username,
		password:   password,
		clientID:   clientID,
		secret:     secret,
		lim:        rl,
	}
	c.BulkService = &BulkService{client: c}
	c.AccountService = &AccountService{client: c}
	c.ContactService = &ContactService{client: c}
	c.OpportunityService = &OpportunityService{client: c}
	c.UserService = &UserService{client: c}
	return c, nil
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Float64 is a helper routine that allocates a new Float64 value
// to store v and returns a pointer to it.
func Float64(v float64) *float64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// makeRequest provides a single function to add common items to the request.
func (c *Client) makeRequest(ctx context.Context, req *http.Request, v interface{}) error {
	token, err := c.getToken()
	if err != nil {
		return fmt.Errorf("error getting token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Accept", "application/json")
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if !c.lim.Allow() {
		c.lim.Wait(ctx)
	}

	rc := req.WithContext(ctx)
	res, err := c.HTTPClient.Do(rc)
	if err != nil {
		return fmt.Errorf("error with do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {

		var salesforceErr error

		switch res.StatusCode {
		case 300:
			salesforceErr = ErrMultipleExternalIDMatch
		case 304:
			salesforceErr = ErrRequestContentNotChanged
		case 400:
			salesforceErr = ErrBadRequest
			var sfbre []BadRequestError
			if err = json.NewDecoder(res.Body).Decode(&sfbre); err == nil {
				fields := strings.Join(sfbre[0].Fields, ",")
				salesforceErr = fmt.Errorf("%w: %s %s", salesforceErr, sfbre[0].Message, fields)
			}
		case 401:
			salesforceErr = ErrUnauthorized
		case 403:
			salesforceErr = ErrForbidden
		case 405:
			salesforceErr = ErrMethodNotAllowed
		case 409:
			salesforceErr = ErrConflict
		case 500:
			salesforceErr = ErrInternalError
		default:
			salesforceErr = ErrUnknown
		}

		var sfbre []BadRequestError
		if err = json.NewDecoder(res.Body).Decode(&sfbre); err == nil {
			fields := strings.Join(sfbre[0].Fields, ",")
			salesforceErr = fmt.Errorf("%w: %s %s", salesforceErr, sfbre[0].Message, fields)
		}

		return salesforceErr

	}

	if res.StatusCode == http.StatusCreated {
		return nil
	}

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}

func (c *Client) getToken() (*sftoken, error) {
	u := fmt.Sprintf("%s/services/oauth2/token", c.BaseURL)
	method := "POST"
	username := url.QueryEscape(c.username)
	password := url.QueryEscape(c.password)
	clientID := url.QueryEscape(c.clientID)
	clientSecret := url.QueryEscape(c.secret)
	pl := fmt.Sprintf("grant_type=password&username=%s&password=%s&client_id=%s&client_secret=%s", username, password, clientID, clientSecret)
	payload := strings.NewReader(pl)

	req, err := http.NewRequest(method, u, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var t sftoken
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return nil, err
	}
	c.token = &t
	return &t, nil
}
