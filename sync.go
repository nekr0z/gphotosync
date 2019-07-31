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

	photoslibrary "github.com/nekr0z/gphotoslibrary"
	"google.golang.org/api/googleapi"
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

	mediaPath := path.Join(lib.Path, strconv.Itoa(remoteCreationTime.Year()), fmt.Sprintf("%02d", remoteCreationTime.Month()), mItem.Filename)

	if err := os.MkdirAll(path.Dir(mediaPath), 0755); err != nil {
		return err
	}

	// multiple items with same filename can exist in remote
	// we try to mitigate it by adding timestamp to local filename
	if stat, err := os.Stat(mediaPath); os.IsNotExist(err) == false && stat.ModTime().UnixNano() != remoteCreationTime.UnixNano() {
		mediaPath = mediaPath + "-" + strconv.FormatInt(remoteCreationTime.UnixNano(), 16) + path.Ext(mediaPath)
	}

	if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
		mediaURL := getDownloadUrl(mItem)

		fmt.Printf("downloading \"%s\" to \"%s\"", mediaURL, mediaPath)
		err := DownloadFile(mediaURL, mediaPath)
		if err != nil {
			return err
		}

		err = os.Chtimes(mediaPath, remoteCreationTime, remoteCreationTime)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("\"%s\" exists", mediaPath)
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
	waitTime := 1 * time.Minute

	for {
		res, err := mediaItemsService.Search(&photoslibrary.SearchMediaItemsRequest{PageToken: pageToken, PageSize: 100}).Do()
		if err != nil {
			if apiError, ok := err.(*googleapi.Error); ok {
				/// If the quota of requests to the Library API is exceeded, the API returns an error code 429 and a message that the project has exceeded the quota.
				if apiError.Code == 429 {
					log.Printf("failed to get media items: %s", apiError.Message)
					time.Sleep(waitTime)
					waitTime = waitTime * 2
					continue
				}
				return err
			}
		}

		fmt.Printf("processing %d items", len(res.MediaItems))
		for _, mItem := range res.MediaItems {
			err = lib.SyncMediaItem(mItem)
			if err != nil {
				log.Printf("failed to sync item \"%s\": %s", mItem.Id, err)
			}
		}

		// if NextPageToken is empty, we reached the last page of items list
		if res.NextPageToken == "" {
			fmt.Printf("syncing is done")
			return nil
		} else {
			fmt.Printf("requesting next page")
		}

		pageToken = res.NextPageToken
		waitTime = 1 * time.Minute
	}
}
