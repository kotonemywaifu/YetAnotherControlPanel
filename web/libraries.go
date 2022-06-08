package web

import (
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

var libraryRegistry = make(map[string]*Library)

type Library struct {
	Name        string
	Type        string
	Version     string
	DownloadUrl string
	VerifyMd5   string
	File        string // generated dynamically
}

var LibraryDir string

var libraryHeader string

func InitializeLibraries(r *gin.Engine) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	LibraryDir = cacheDir + "/yacp/libraries"

	// mkdir if not exists
	if _, err := os.Stat(LibraryDir); os.IsNotExist(err) {
		err = os.MkdirAll(LibraryDir, 0755)
		if err != nil {
			return err
		}
	}

	// add libraries to registry
	libraryRegistry["swal"] = &Library{
		Name:        "sweetalert",
		Type:        "js",
		Version:     "11.4.17",
		DownloadUrl: "https://cdn.jsdelivr.net/npm/sweetalert2@11.4.17/dist/sweetalert2.all.min.js",
		VerifyMd5:   "eeda8bbd18f70c496ab3c2d29667d1fe",
	}
	libraryRegistry["jquery"] = &Library{
		Name:        "jquery",
		Type:        "js",
		Version:     "3.6.0",
		DownloadUrl: "https://cdn.jsdelivr.net/npm/jquery@3.6.0/dist/jquery.min.js",
		VerifyMd5:   "8fb8fee4fcc3cc86ff6c724154c49c42",
	}
	libraryRegistry["tailwindcss"] = &Library{
		Name:        "tailwindcss",
		Type:        "css",
		Version:     "2.2",
		DownloadUrl: "https://cdn.jsdelivr.net/npm/tailwindcss@2.2/dist/tailwind.min.css",
		VerifyMd5:   "e35af4d8ceb624072098fa9a3d970aaa",
	}

	// download libraries
	for _, library := range libraryRegistry {
		libn := library.Name + "-" + library.Version + "." + library.Type
		library.File = LibraryDir + "/" + libn
		err = PrepareLibrary(library)
		if err != nil {
			return errors.New("failed to download library " + library.Name + " :" + err.Error())
		}
		if library.Type == "js" {
			libraryHeader += "<script src=\"/assets/" + libn + "\"></script>\n"
		} else if library.Type == "css" {
			libraryHeader += "<link rel=\"stylesheet\" href=\"/assets/" + libn + "\"></link>\n"
		}
	}

	// add routes
	r.Static("/assets", LibraryDir)

	return nil
}

func PrepareLibrary(library *Library) error {
	// download library if not exists
	isFreshDownload := false
	if _, err := os.Stat(library.File); os.IsNotExist(err) {
		log.Println("Downloading library " + library.Name + " to " + library.File)
		err := util.DownloadFile(library.DownloadUrl, library.File)
		if err != nil {
			return err
		}
		isFreshDownload = true
	}

	// verify md5
	if library.VerifyMd5 != "" {
		verified, err := util.VerifyMd5(library.File, library.VerifyMd5)
		if err != nil {
			return err
		}
		if !verified {
			if isFreshDownload {
				return errors.New("md5 verification failed for library " + library.Name)
			} else {
				log.Println("md5 verification failed for library " + library.Name + ", downloading again")
				// delete file and download again
				err = os.Remove(library.File)
				if err != nil {
					return err
				}
				return PrepareLibrary(library)
			}
		}
	}

	return nil
}

/**
 * Try push library through http 2
 * @return library in html header
 */
func tryPushLibrary(c *gin.Context) string {
	// resourcePushed := false
	// if pusher := c.Writer.Pusher(); pusher != nil {
	// 	if err := pusher.Push("", nil)
	// }
	// TODO: implement
	return libraryHeader
}
