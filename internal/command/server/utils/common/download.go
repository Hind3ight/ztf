package serverUtils

import (
	"archive/zip"
	"fmt"
	i118Utils "github.com/aaronchen2k/deeptest/internal/command/utils/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	errUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/err"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Download(uri string, dst string) error {
	if strings.Index(uri, "?") < 0 {
		uri += "?"
	} else {
		uri += "&"
	}
	uri += fmt.Sprintf("&r=%d", time.Now().Unix())

	res, err := http.Get(uri)
	if err != nil {
		logUtils.PrintTo(i118Utils.Sprintf("download_fail", uri, err.Error()))
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logUtils.PrintTo(i118Utils.Sprintf("download_read_fail", uri, err.Error()))
	}

	err = ioutil.WriteFile(dst, bytes, 0666)
	if err != nil {
		logUtils.PrintTo(i118Utils.Sprintf("download_write_fail", dst, err.Error()))
	} else {
		logUtils.PrintTo(i118Utils.Sprintf("download_success", uri, dst))
	}

	return err
}

func GetZipSingleDir(path string) string {
	folder := ""
	z := archiver.Zip{}
	err := z.Walk(path, func(f archiver.File) error {
		if f.IsDir() {
			zfh, ok := f.Header.(zip.FileHeader)
			if ok {
				//logUtils.PrintTo("file: " + zfh.Name)

				if folder == "" && zfh.Name != "__MACOSX" {
					folder = zfh.Name
				} else {
					if strings.Index(zfh.Name, folder) != 0 {
						return errUtils.New("found more than one folder")
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return ""
	}

	return folder
}
