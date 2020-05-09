package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	mediaIdsPath = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "tmp/mediaIds.json")
	backupPath   = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "tmp/bkp")
)

type mediaInfo struct {
	Id  string `json:"id"`
	URL string `json:"url"`
}

func (cli *CLI) listMediaItems(url string) *MediaItemsListResponse {

	resp, err := cli.client.Get(url)
	if err != nil {
		log.Fatalf("could not fetch the media items. Error %v", err)
	}
	mi := &MediaItemsListResponse{}
	json.NewDecoder(resp.Body).Decode(mi)
	return mi
}

func (cli *CLI) getJunk(url string) {

	mediaItems := make([]*mediaInfo, 0)
	mi := cli.listMediaItems(url)

	for _, v := range mi.MediaItems {
		width, _ := strconv.Atoi(v.MediaMetadata.Width)
		height, _ := strconv.Atoi(v.MediaMetadata.Height)
		size := calculateImageSize(width, height)
		if size < 500 && v.MediaMetadata.Photo != nil && v.MediaMetadata.Photo.FocalLength == 0 && v.MediaMetadata.Photo.CameraModel == "" {
			item := &mediaInfo{
				Id:  v.ID,
				URL: v.BaseURL,
			}
			mediaItems = append(mediaItems, item)
		}

	}

	_, err := os.Stat(mediaIdsPath)
	if err == nil {
		items := cli.getMediaIds()
		mediaItems = append(mediaItems, items...)

	}

	f, err := os.OpenFile(mediaIdsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		log.Fatalf("Unable to cache the mediaitems Ids")
	}
	defer f.Close()

	b, err := json.Marshal(mediaItems)
	if err != nil {
		log.Fatalf("unable to unmarshall. Err %v", err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Fatalf("unable to write data. Err %v", err)
	}
	if mi.NextPageToken != "" {
		cli.getJunk(listMediaItems + "&pageToken=" + mi.NextPageToken)
	}

}

func (cli *CLI) getMediaIds() []*mediaInfo {
	info := make([]*mediaInfo, 0)
	b, err := ioutil.ReadFile(mediaIdsPath)
	if err != nil {
		log.Fatal("Could not read the file:", mediaIdsPath)
	}
	json.Unmarshal(b, &info)
	return info
}

func calculateImageSize(width, height int) int {
	return ((width * height * 16) / 8) / 1024
}

func (cli *CLI) createBackup(list []*mediaInfo) {

	var wg sync.WaitGroup
	err := os.MkdirAll(backupPath, 0700)
	if err != nil {
		log.Fatalf("Could not create the backup directory. Error %v", err)
	}

	wg.Add(len(list))
	for i, v := range list {
		go func(index int, mi *mediaInfo) {
			defer wg.Done()
			fn := fmt.Sprintf("%s/%d", backupPath, index)
			f, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0700)
			if err != nil {
				log.Println("Couldn't create file at ", backupPath, index, err)
				return
			}
			defer f.Close()
			r, err := cli.client.Get(mi.URL)
			if err != nil {
				log.Println("Couldn't get file from url ", mi.URL, err)
				return
			}
			defer r.Body.Close()
			_, err = io.Copy(f, r.Body)
			if err != nil {
				log.Println("could not create download file to ", backupPath, index, err)
				return
			}
		}(i, v)
	}

	wg.Wait()
}

func (cli *CLI) RemoveMediaItems(albumId string, list []*mediaInfo) {
	arr := splitItems(50, list)

	for i := 0; i < len(arr); i++ {
		ids := make([]string, 0)
		for j := 0; j < len(arr[i]); j++ {
			ids = append(ids, arr[i][j].Id)
		}
		cli.BatchDelete(albumId, ids)
	}

}

func splitItems(limit int, list []*mediaInfo) [][]*mediaInfo {

	var res [][]*mediaInfo

	for i := 0; i < len(list); i += limit {
		end := i + limit

		if end > len(list) {
			end = len(list)
		}

		res = append(res, list[i:end])
	}
	return res
}

func (cli *CLI) BatchDelete(albumId string, ids []string) {

	type req struct {
		MediaItemIds []string `json:"mediaItemIds"`
	}

	reqObj := &req{
		MediaItemIds: ids,
	}

	url := fmt.Sprintf("https://photoslibrary.googleapis.com/v1/albums/%s:batchRemoveMediaItems", albumId)

	rb, err := json.Marshal(reqObj)
	if err != nil {
		log.Fatalf("unable to marhsal array of ids. Err: %v", err)
	}
	resp, err := cli.client.Post(url, "application/json", bytes.NewBuffer(rb))
	if err != nil {
		log.Printf("Deletion failed for a request. Err %v", err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("\n%s\n", b)

	fmt.Println("Status: ", resp.Status)
}
