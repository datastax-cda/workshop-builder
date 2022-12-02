package util

import (
	"compress/flate"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/mholt/archiver"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func CloneRepo(repoPath string, destinationPath string) error {
	data, err := base64.StdEncoding.DecodeString("Z2hwX3Q0bk83aWg3a3lUTTV1aWRyQXdwb0x4OTM1T09JdzBqRUV1Mwo=")
	if err != nil {
		return err
	}

	gitToken := strings.TrimSpace(string(data))

	fmt.Println("git clone " + repoPath)

	_, err = git.PlainClone(
		destinationPath,
		false,
		&git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "doesnotmatter",
				Password: gitToken,
			},
			URL:      repoPath,
			Progress: os.Stdout,
		},
	)

	if err != nil {
		return fmt.Errorf("cannot clone base git repo + %+v", err)
	}
	return nil
}

func RemoveGitMetadata(destinationPath string) error {

	err := os.RemoveAll(destinationPath + "/.gitignore")
	if err != nil {
		fmt.Println("could not remove .gitignore file")
	}
	err = os.RemoveAll(destinationPath + "/.git")
	if err != nil {
		fmt.Println("could not remove .git directory")
	}
	return nil
}

func ZipIt(source, target string) error {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}

	err := z.Archive([]string{source}, target)
	if err != nil {
		return err
	}

	return nil

}
