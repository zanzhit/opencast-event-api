package parser

import (
	"os/exec"
	"strings"
)

func videoDuration(fileNameWithPath string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", fileNameWithPath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		output := string(out)

		if strings.Contains(output, "Duration") {
			start := strings.Index(output, "Duration") + 10
			end := start + 8

			return output[start:end], nil
		}
	}

	return "", err
}
