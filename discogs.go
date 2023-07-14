package discogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	discogsAPI = "https://api.discogs.com"
)

// Options is a set of options to use discogs API client
type Options struct {
	// Discogs API endpoint (optional).
	URL string
	// Currency to use (optional, default is USD).
	Currency string
	// UserAgent to to call discogs api with.
	UserAgent string
	// Token provided by discogs (optional).
	Token string
	// HTTP client instance to use for HTTP requests
	Client *http.Client
	// Rate limit instance to track request rates
	RateLimit *RateLimit
}

// Discogs is an interface for making Discogs API requests.
type Discogs interface {
	CollectionService
	DatabaseService
	MarketPlaceService
	SearchService
	WantlistService
}

type discogs struct {
	CollectionService
	DatabaseService
	SearchService
	MarketPlaceService
	WantlistService
}

type requestFunc func(ctx context.Context, path string, params url.Values, resp interface{}) error
type writeFunc func(ctx context.Context, path string, method string, params url.Values, payload interface{}, resp interface{}, successStatus int) error

// New returns a new discogs API client.
func New(o *Options) (Discogs, error) {
	header := &http.Header{}

	if o == nil || o.UserAgent == "" {
		return nil, ErrUserAgentInvalid
	}

	header.Add("User-Agent", o.UserAgent)

	cur, err := currency(o.Currency)
	if err != nil {
		return nil, err
	}

	// set token, it's required for some queries like search
	if o.Token != "" {
		header.Add("Authorization", "Discogs token="+o.Token)
	}

	if o.URL == "" {
		o.URL = discogsAPI
	}

	client := o.Client
	if client == nil {
		client = &http.Client{}
	}
	req := func(ctx context.Context, path string, params url.Values, resp interface{}) error {
		return request(ctx, client, "GET", header, o.RateLimit, path, params, nil, resp, http.StatusOK)
	}

	write := func(ctx context.Context, path string, method string, params url.Values, payload interface{}, resp interface{}, successStatus int) error {
		return request(ctx, client, method, header, o.RateLimit, path, params, payload, resp, successStatus)
	}

	impl := discogs{
		newCollectionService(req, o.URL+"/users"),
		newDatabaseService(req, o.URL, cur),
		newSearchService(req, o.URL+"/database/search"),
		newMarketPlaceService(req, o.URL+"/marketplace", cur),
		newWantlistService(req, write, o.URL+"/users"),
	}

	if o.RateLimit != nil {
		return o.RateLimit.Client(impl), nil
	}

	return impl, nil
}

// currency validates currency for marketplace data.
// Defaults to the authenticated users currency. Must be one of the following:
// USD GBP EUR CAD AUD JPY CHF MXN BRL NZD SEK ZAR
func currency(c string) (string, error) {
	switch c {
	case "USD", "GBP", "EUR", "CAD", "AUD", "JPY", "CHF", "MXN", "BRL", "NZD", "SEK", "ZAR":
		return c, nil
	case "":
		return "USD", nil
	default:
		return "", ErrCurrencyNotSupported
	}
}

func request(ctx context.Context, client *http.Client, method string, header *http.Header,
	rl *RateLimit, path string, params url.Values, payload interface{}, resp interface{}, successStatus int) error {

	var rawPayload io.Reader
	if payload != nil {
		s, _ := json.Marshal(payload)
		rawPayload = bytes.NewReader(s)
	}

	r, err := http.NewRequestWithContext(ctx, method, path+"?"+params.Encode(), rawPayload)
	if err != nil {
		return err
	}
	r.Header = *header
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if rl != nil {
		total, _ := strconv.Atoi(response.Header.Get("X-Discogs-Ratelimit"))               // The total number of requests you can make in a one minute window.
		used, _ := strconv.Atoi(response.Header.Get("X-Discogs-Ratelimit-Used"))           // The number of requests you’ve made in your existing rate limit window.
		remaining, _ := strconv.Atoi(response.Header.Get("X-Discogs-Ratelimit-Remaining")) // The number of remaining requests you are able to make in the existing rate limit window.
		rl.Update(total, used, remaining)
	}

	if response.StatusCode != successStatus {
		switch response.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusTooManyRequests:
			return ErrTooManyRequests
		default:
			return fmt.Errorf("unknown error: %s", response.Status)
		}
	}

	if resp != nil {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return json.Unmarshal(body, &resp)
	}

	return nil
}
