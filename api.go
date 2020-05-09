package main

const (
	listAlbums     = "https://photoslibrary.googleapis.com/v1/albums?pageSize=50"
	listMediaItems = "https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=100"
)

type Album struct {
	ID                    string    `json:"id"`
	Title                 string    `json:"title"`
	ProductURL            string    `json:"productUrl"`
	IsWriteable           string    `json:"isWriteable"`
	MediaItemsCount       string    `json:"mediaItemsCount"`
	CoverPhotoBaseURL     string    `json:"coverPhotoBaseUrl"`
	CoverPhotoMediaItemID string    `json:"coverPhotoMediaItemId"`
	ShareInfo             ShareInfo `json:"shareInfo"`
}

type ShareInfo struct {
	SharedAlbumOptions SharedAlbumOptions `json:"sharedAlbumOptions"`
	ShareableURL       string             `json:"shareableUrl"`
	ShareToken         string             `json:"shareToken"`
	IsJoined           string             `json:"isJoined"`
	IsOwned            string             `json:"isOwned"`
}

type SharedAlbumOptions struct {
	IsCollaborative string `json:"isCollaborative"`
	IsCommentable   string `json:"isCommentable"`
}

type AlbumListResponse struct {
	Albums        []Album `json:"albums"`
	NextPageToken string  `json:"nextPageToken"`
}

//MediaItems

type Photo struct {
	CameraMake      string `json:"cameraMake"`
	CameraModel     string `json:"cameraModel"`
	FocalLength     int    `json:"focalLength"`
	ApertureFNumber int    `json:"apertureFNumber"`
	IsoEquivalent   int    `json:"isoEquivalent"`
	ExposureTime    string `json:"exposureTime"`
}

type Video struct {
	CameraMake  string `json:"cameraMake"`
	CameraModel string `json:"cameraModel"`
	Fps         int    `json:"fps"`
}

type ContributorInfo struct {
	ProfilePictureBaseURL string `json:"profilePictureBaseUrl"`
	DisplayName           string `json:"displayName"`
}

type MediaMetadata struct {
	CreationTime string `json:"creationTime"`
	Width        string `json:"width"`
	Height       string `json:"height"`
	Photo        *Photo `json:"photo"`
	Video        *Video `json:"video"`
}

type MediaItem struct {
	ID              string           `json:"id"`
	Description     string           `json:"description"`
	ProductURL      string           `json:"productUrl"`
	BaseURL         string           `json:"baseUrl"`
	MimeType        string           `json:"mimeType"`
	MediaMetadata   *MediaMetadata   `json:"mediaMetadata"`
	ContributorInfo *ContributorInfo `json:"contributorInfo"`
	Filename        string           `json:"filename"`
}

type MediaItemsListResponse struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}
