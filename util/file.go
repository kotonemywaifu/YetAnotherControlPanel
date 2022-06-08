package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var HttpClient = &http.Client{
	Timeout: time.Second * 10,
	CheckRedirect: func(r *http.Request, via []*http.Request) error {
		r.URL.Opaque = r.URL.Path
		return nil
	},
}

func DownloadFile(url, file string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:100.0) Gecko/20100101 Firefox/100.0")
	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func VerifyMd5(file, md5Hash string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	hash := hex.EncodeToString(h.Sum(nil))
	return hash == md5Hash, nil
}

func ReadFile(fs fs.FS, file string) (string, error) {
	f, err := fs.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	result, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func Must(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}
