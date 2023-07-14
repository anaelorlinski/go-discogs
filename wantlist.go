package discogs

import (
	"context"
	"fmt"
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
	url     string
}

func newWantlistService(req requestFunc, url string) WantlistService {
	return &wantlistService{
		request: req,
		url:     url,
	}
}

// CollectionFolders serves collection response from discogs.
type WantlistItem struct {
	ID          int              `json:"id"`
	Rating      int              `json:"rating"`
	Notes       string           `json:"notes"`
	NotesPublic string           `json:"notes_public"`
	Release     BasicInformation `json:"basic_information"`
}

type WantlistPagination struct {
	Page    int
	PerPage int
}

type Wantlist struct {
	Pagination Page           `json:"pagination"`
	Items      []WantlistItem `json:"wants"`
}

func (w *Wantlist) Merge(other *Wantlist) {
	w.Items = append(w.Items, other.Items...)
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

func (s *wantlistService) AddWantlistItem(ctx context.Context, username string, item WantlistItem) error {
	return nil
}

func (s *wantlistService) UpdateWantlistItem(ctx context.Context, username string, item WantlistItem) error {
	return nil
}

func (s *wantlistService) DeleteWantlistItem(ctx context.Context, username string, item WantlistItem) error {
	return nil
}
