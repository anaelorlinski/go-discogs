package discogs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const maxWishlistPageItems = 100

// CollectionService is an interface to work with collection.
type WantlistService interface {
	Wantlist(ctx context.Context, username string, pagination *WantlistPagination) (*Wantlist, error)
	CompleteWantlist(ctx context.Context, username string) (*Wantlist, error)
	AddWantlistItem(ctx context.Context, username string, item WantlistItem) error
	UpdateWantlistItem(ctx context.Context, username string, item WantlistItem) error
	DeleteWantlistItem(ctx context.Context, username string, item WantlistItem) error
}

type wantlistService struct {
	request requestFunc
	write   writeFunc
	url     string
}

func newWantlistService(req requestFunc, write writeFunc, url string) WantlistService {
	return &wantlistService{
		request: req,
		write:   write,
		url:     url,
	}
}

// note : rating is not implemented
type WantlistItem struct {
	ID      int               `json:"id"`
	Notes   string            `json:"notes"`
	Release *BasicInformation `json:"basic_information,omitempty"`
}

type Wantlist struct {
	Pagination Page           `json:"pagination"`
	Items      []WantlistItem `json:"wants"`
}

func (w *Wantlist) Merge(other *Wantlist) {
	w.Items = append(w.Items, other.Items...)
}

type WantlistPagination struct {
	Page    int
	PerPage int
}

// toParams converts pagaination params to request values
func (p *WantlistPagination) params() url.Values {
	if p == nil {
		return nil
	}

	params := url.Values{}
	params.Set("page", strconv.Itoa(p.Page))
	params.Set("per_page", strconv.Itoa(p.PerPage))
	return params
}

func (s *wantlistService) Wantlist(ctx context.Context, username string, pagination *WantlistPagination) (*Wantlist, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}
	var wantlistPage *Wantlist
	err := s.request(ctx, s.url+"/"+username+"/wants", pagination.params(), &wantlistPage)
	return wantlistPage, err

}

func (s *wantlistService) CompleteWantlist(ctx context.Context, username string) (*Wantlist, error) {
	// getting the complete wantlist will trigger many requests
	// so it needs the client to enforce rate limits
	// implementation is done in ratelimited client
	// pass a ratelimite option to the client to create it
	return nil, fmt.Errorf("implemented only on ratelimited")
}

func (s *wantlistService) AddWantlistItem(ctx context.Context, username string, item WantlistItem) (e error) {
	if username == "" {
		return ErrInvalidUsername
	}

	uri := s.url + "/" + username + "/wants/" + strconv.Itoa(item.ID)
	e = s.write(ctx, uri, "PUT", url.Values{}, item, nil, http.StatusCreated)

	return
}

func (s *wantlistService) UpdateWantlistItem(ctx context.Context, username string, item WantlistItem) (e error) {
	if username == "" {
		return ErrInvalidUsername
	}

	uri := s.url + "/" + username + "/wants/" + strconv.Itoa(item.ID)
	e = s.write(ctx, uri, "POST", url.Values{}, item, nil, http.StatusOK)

	return
}

func (s *wantlistService) DeleteWantlistItem(ctx context.Context, username string, item WantlistItem) (e error) {
	if username == "" {
		return ErrInvalidUsername
	}
	uri := s.url + "/" + username + "/wants/" + strconv.Itoa(item.ID)
	e = s.write(ctx, uri, "DELETE", url.Values{}, nil, nil, http.StatusNoContent)

	return
}
