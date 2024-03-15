package main

import (
	"fmt"
	"log"
	"time"

	"opencast/pkg/api/request"
	"opencast/pkg/parser"
)

func main() {
	config, cameras, err := parser.ReadConfig("../config/cameras.json", "../config/config.json")
	if err != nil {
		log.Println("error opening cameras and config:", err)
		return
	}

	acl, processing, err := parser.ReadACLProcessing("../config/acl.json", "../config/process.json")
	if err != nil {
		log.Println("error opening acl and process:", err)
		return
	}

	metadatas, presenters, shinobiFiles, err := parser.ParseCameras("//wsl$/Ubuntu/home/Shinobi/videos/vZLeZUNbaG", cameras)
	if err != nil {
		log.Println("error parse cameras:", err)
		return
	}

	for i := 0; i < len(metadatas); i++ {
		respBody, err := request.SendMultipartRequest(presenters[i], acl, metadatas[i], processing, config)
		if err != nil {
			log.Println("error sending request", err)
			return
		}

		fmt.Println(string(respBody))

		time.Sleep(time.Second * 1)

		deleteRespStatus, err := request.DeleteVideo(config, shinobiFiles[i])
		if err != nil {
			log.Println("error deleting Shinobi video", err)
		}

		fmt.Println(deleteRespStatus)
	}
}
