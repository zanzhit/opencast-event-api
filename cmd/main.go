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

	camerasData, err := parser.ParseCameras("//wsl$/Ubuntu/home/Shinobi/videos/vZLeZUNbaG", cameras)
	if err != nil {
		log.Println("error parse cameras:", err)
		return
	}

	for i := 0; i < len(camerasData); i++ {
		respBody, err := request.SendMultipartRequest(camerasData[i], acl, processing, config)
		if err != nil {
			log.Println("error sending request", err)
			return
		}

		fmt.Println(string(respBody))

		time.Sleep(time.Second * 1)

		deleteRespStatus, err := request.DeleteVideo(config, camerasData[i].ShinobiFile)
		if err != nil {
			log.Println("error deleting Shinobi video", err)
		}

		fmt.Println(deleteRespStatus)
	}
}
