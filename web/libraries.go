package web

import (
	"errors"
	"log"
	"os"

	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

var libraryRegistry = make(map[string]*Library)

type Library struct {
	Name        string
	Type        string
	Version     string
	DownloadUrl string
	VerifyMd5   string
}

var LibraryDir string

func InitializeLibraries() error {
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

	// download libraries
	for _, library := range libraryRegistry {
		err = PrepareLibrary(library)
		if err != nil {
			return errors.New("failed to download library " + library.Name + " :" + err.Error())
		}
	}

	return nil
}

func PrepareLibrary(library *Library) error {
	libraryFile := LibraryDir + "/" + library.Name + "-" + library.Version + "." + library.Type

	// download library if not exists
	isFreshDownload := false
	if _, err := os.Stat(libraryFile); os.IsNotExist(err) {
		log.Println("Downloading library " + library.Name + " to " + libraryFile)
		err := util.DownloadFile(library.DownloadUrl, libraryFile)
		if err != nil {
			return err
		}
		isFreshDownload = true
	}

	// verify md5
	if library.VerifyMd5 != "" {
		verified, err := util.VerifyMd5(libraryFile, library.VerifyMd5)
		if err != nil {
			return err
		}
		if !verified {
			if isFreshDownload {
				return errors.New("md5 verification failed for library " + library.Name)
			} else {
				log.Println("md5 verification failed for library " + library.Name + ", downloading again")
				// delete file and download again
				err = os.Remove(libraryFile)
				if err != nil {
					return err
				}
				return DownloadLibrary(library)
			}
		}
	}

	return nil
}
