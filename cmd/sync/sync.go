package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"google.golang.org/api/googleapi"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

func DownloadFile(url string, path string) error {

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getDownloadUrl(mediaItem *photoslibrary.MediaItem) string {
	if mediaItem.MediaMetadata.Photo == nil {
		return fmt.Sprintf("%s=dv", mediaItem.BaseUrl)
	}

	return fmt.Sprintf("%s=d", mediaItem.BaseUrl)

}

func (lib *Library) SyncMediaItem(mItem *photoslibrary.MediaItem) error {
	remoteCreationTime, err := time.Parse(time.RFC3339, mItem.MediaMetadata.CreationTime)
	if err != nil {
		return err
	}

	mediaPath := path.Join(lib.Path, strconv.Itoa(remoteCreationTime.Year()), mItem.Filename)

	if err := os.MkdirAll(path.Dir(mediaPath), 0755); err != nil {
		return err
	}

	if stat, err := os.Stat(mediaPath); os.IsNotExist(err) || (os.IsExist(err) && stat.ModTime().UnixNano() <= remoteCreationTime.UnixNano()) {
		mediaURL := getDownloadUrl(mItem)

		log.Printf("downloading \"%s\" to \"%s\"", mediaURL, mediaPath)
		err := DownloadFile(mediaURL, mediaPath)
		if err != nil {
			return err
		}

		err = os.Chtimes(mediaPath, remoteCreationTime, remoteCreationTime)
		if err != nil {
			return err
		}
	} else {
		log.Printf("\"%s\" exists", mediaPath)
	}

	return nil
}

func (lib *Library) GetTokenPath() string {
	return path.Join(lib.Path, ".token.json")
}

func (lib *Library) Sync() error {
	ctx := context.Background()

	oauthClient, err := NewOAuthClient(ctx, GoogleClientId, GoogleClientSecret, lib.GetTokenPath())
	if err != nil {
		return err
	}
	photoslibraryService, err := photoslibrary.New(oauthClient)
	if err != nil {
		return err
	}

	mediaItemsService := photoslibrary.NewMediaItemsService(photoslibraryService)
	pageToken := ""

	for {
		res, err := mediaItemsService.Search(&photoslibrary.SearchMediaItemsRequest{PageToken: pageToken, PageSize: 100}).Do()
		if err != nil {
			if apiError, ok := err.(*googleapi.Error); ok {
				/// If the quota of requests to the Library API is exceeded, the API returns an error code 429 and a message that the project has exceeded the quota.
				if apiError.Code == 429 {
					log.Printf("failed to get media items: %s", apiError.Message)
					time.Sleep(5 * time.Minute)
					continue
				}
				return err
			}
		}

		// no more photos there
		if len(res.MediaItems) == 0 {
			log.Printf("syncing is done")
			return nil
		} else {
			log.Printf("processing %d items", len(res.MediaItems))
		}

		for _, mItem := range res.MediaItems {
			err = lib.SyncMediaItem(mItem)
			if err != nil {
				log.Printf("failed to sync item \"%s\": %s", mItem.Id, err)
			}
		}

		pageToken = res.NextPageToken
	}
}
