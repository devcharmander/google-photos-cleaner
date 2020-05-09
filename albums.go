package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var count = 0

func (cli *CLI) ListAlbums(url string) (int, map[int]string) {
	resp, err := cli.client.Get(url)
	if err != nil {
		log.Fatalf("error getting list of albums. Err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Status code: ", resp.StatusCode)
		bd, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", bd)
		fmt.Println(resp.Status)
	}
	album := &AlbumListResponse{}
	err = json.NewDecoder(resp.Body).Decode(album)
	if err != nil {
		log.Fatalf("error while decoding response. Err:%v", err)
	}

	akv := make(map[int]string)
	for i, v := range album.Albums {
		//akv[i] = v.Title + "##" + v.MediaItemsCount // didn't want to create a new struct here
		akv[i] = v.ID
		fmt.Println(i+1, ": ", v.Title)
	}
	count += len(album.Albums)
	if album.NextPageToken != "" {
		cli.ListAlbums(listAlbums + "&pageToken=" + album.NextPageToken)
	}
	return count, akv
}
