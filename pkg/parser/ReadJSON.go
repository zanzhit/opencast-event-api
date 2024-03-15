package parser

import "os"

func ReadACLProcessing(aclJSON, processingJSON string) (acl []byte, processing []byte, err error) {

	acl, err = os.ReadFile(aclJSON)
	if err != nil {
		return nil, nil, err
	}

	processing, err = os.ReadFile(processingJSON)
	if err != nil {
		return nil, nil, err
	}

	return acl, processing, nil
}

func ReadConfig(camerasJSON, configJSON string) (config []byte, cameras []byte, err error) {
	config, err = os.ReadFile(configJSON)
	if err != nil {
		return nil, nil, err
	}

	cameras, err = os.ReadFile(camerasJSON)
	if err != nil {
		return nil, nil, err
	}

	return config, cameras, nil
}
