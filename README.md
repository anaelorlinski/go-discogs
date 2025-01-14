# REST API 2.0 Discogs.com client

go-discogs is a Go client library for the [Discogs API](https://www.discogs.com/developers/). Check the usage section to see how to access the Discogs API.

Note : this is a fork of sblinch/go-discogs itself a fork of irlndts/go-discogs

The lib is under MIT but be sure you are familiar with [Discogs API Terms of Use](https://support.discogs.com/hc/en-us/articles/360009334593-API-Terms-of-Use).

### Features
 * Database
    * [Releases](#releases)
    * Release Rating
    * Master Releases
    * Master Versions
    * Artists
    * Artist Releases
    * Label
    * All Label Releases
 * [Search](#search)
 * [User Collection](#user-collection)
    * Collection Folders
    * Folder
    * Collection Items by Folder
    * Collection Items by Release
 * [Marketplace](#marketplace)
    * Price Suggestions
    * Release Statistics
 
Install
--------
    go get github.com/anaelorlinski/go-discogs

Usage
---------
The discogs package provides a client for accessing the Discogs API. 
First of all import library and init client variable. According to discogs api documentation you [must provide your user-agent](https://www.discogs.com/developers/#page:home,header:home-general-information). 
```go
import "github.com/anaelorlinski/go-discogs"
```

Some requests require authentication (as any user). According to [Discogs](https://www.discogs.com/developers/#page:authentication,header:authentication-discogs-auth-flow), to send requests with Discogs Auth, you have two options: sending your credentials in the query string with key and secret parameters or a [token parameter](https://www.discogs.com/settings/developers).

```go
client, err := discogs.New(&discogs.Options{
        UserAgent: "Some Name",
        Currency:  "EUR", // optional, "USD" (default), "GBP", "EUR", "CAD", "AUD", "JPY", "CHF", "MXN", "BRL", "NZD", "SEK", "ZAR" are allowed
        Token:     "Some Token", // optional
        URL:       "https://api.discogs.com", // optional
    })
``` 

#### Releases
```go
  release, _ := client.Release(context.Background(), 9893847)
  fmt.Println(release.Artists[0].Name, " - ", release.Title) 
  // St. Petersburg Ska-Jazz Review  -  Elephant Riddim
```

#### Search
Issue a search query to discogs database. This endpoint accepts pagination parameters.
Authentication (as any user) is required.

Use `SearchRequest` struct to create a request.
```go
type SearchRequest struct {
    Q             string // search query (optional)
    Type          string // one of release, master, artist, label (optional)
    Title         string // search by combined “Artist Name - Release Title” title field (optional)
    ReleaseTitle string // search release titles (optional)
    Credit        string // search release credits (optional)
    Artist        string // search artist names (optional)
    Anv           string // search artist ANV (optional)
    Label         string // search label names (optional)
    Genre         string // search genres (optional)
    Style         string // search styles (optional)
    Country       string // search release country (optional)
    Year          string // search release year (optional)
    Format        string // search formats (optional)
    Catno         string // search catalog number (optional)
    Barcode       string // search barcodes (optional)
    Track         string // search track titles (optional)
    Submitter     string // search submitter username (optional)
    Contributer   string // search contributor usernames (optional)

    Page     int // optional
    PerPage  int // optional
}
```

```go
  request := discogs.SearchRequest{Artist: "reggaenauts", ReleaseTitle: "river rock", Page: 0, PerPage: 1}
  search, _ := client.Search(context.Background(), request)

  for _, r := range search.Results {
    fmt.Println(r.Title)
  }
```

#### User Collection

Query a users [collection](https://www.discogs.com/developers#page:user-collection).

##### Collection Folders
```go
  collection, err := client.CollectionFolders(context.Background(), "my_user")
```
##### Folder
```go
  folder, err := client.Folder(context.Background(), "my_user", 0)
```
##### Collection Items by Folder
```go
  items, err := client.CollectionItemsByFolder(context.Background(), "my_user", 0, &Pagination{Sort: "artist", SortOrder: "desc", PerPage: 2})
```
##### Collection Items by Release
```go
  items, err := client.CollectionItemsByRelease(context.Background(), "my_user", 12934893)
```

#### Marketplace

Query a user's [marketplace](https://www.discogs.com/developers/#page:marketplace)

##### Price Suggestions

Retrieve price suggestions for the provided Release ID

```go
  suggestions, err := client.PriceSuggestions(context.Background(), 12345)
```

##### Release Statistics

Retrieve marketplace statistics for the provided Release ID

```go
  stats, err := client.ReleaseStatistics(context.Background(), 12345)
```

...
