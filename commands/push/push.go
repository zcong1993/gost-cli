package push

import (
	"fmt"
	"github.com/gost-c/gost-cli/commands"
	"github.com/gost-c/gost-cli/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	u "github.com/zcong1993/utils"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"os"
	"path"
)

var url = utils.BaseURL + "api/create"

// Gist is struct of push api's post data
type Gist struct {
	Public      bool
	Description string
	Version     uint
	Files       []File
}

// File is struct of gist's file
type File struct {
	Filename string
	Content  string
}

// Run is sub command runner for push
func Run(files []string, description string) {
	f, err := getFiles(files)
	if err != nil {
		log.Fatalf("parse files error: %v", err)
	}
	gist := Gist{Public: true, Description: description, Files: f, Version: 1}
	token, err := utils.ReadConfig()
	if err != nil || token == nil {
		log.Fatal("Get token failed, please login first")
	}

	var res commands.Result
	err = u.PostJSON(url, gist, &res, map[string]string{"Authorization": "Bearer " + string(token)})
	if err != nil {
		log.Fatalf("Unexpected error occurred: %s, make sure you have logged in", err.Error())
	}

	if res.Code != "200" {
		log.Fatalf("Push error: %s", res.Msg)
	}

	log.Infof("Push success! The url is %s%s", utils.WebURL, res.Msg)
}

func getFiles(files []string) ([]File, error) {
	eg := errgroup.Group{}
	var results []File
	for _, file := range files {
		file := file
		eg.Go(func() error {
			content, err := ioutil.ReadFile(file)
			fmt.Fprintf(os.Stdout, "--> Parsing file: %15s\n", file)
			if err != nil {
				return errors.Wrapf(err,
					"failed to get file content: %s", file)
			}
			results = append(results, File{path.Base(file), string(content)})
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, errors.Wrap(err, "one of the goroutines failed")
	}
	return results, nil
}