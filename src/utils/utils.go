package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

// ReadAnswer prints out a yes/no form with string from `info`
// and returns boolean value based on user input (y/Y or n/N) or
// return `defaultAnswer` if input is omitted.
// If input is neither of them, print form again.
func ReadAnswer(info string, defaultAnswer bool) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(info)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\r", "", 1)
	text = strings.Replace(text, "\n", "", 1)
	if len(text) == 0 {
		return defaultAnswer
	} else if text == "y" || text == "Y" {
		return true
	} else if text == "n" || text == "N" {
		return false
	}
	return ReadAnswer(info, defaultAnswer)
}

// CheckExistAndCreate checks folder existence
// and make that folder if it does not exist
func CheckExistAndCreate(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		os.Mkdir(dir, 0700)
	}
}

// Unzip unzips zip
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0700)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, 0700)
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy .
func Copy(src, dest string, recursive bool, filters []string) error {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	os.MkdirAll(dest, 0700)

	for _, file := range dir {
		fSrcPath := filepath.Join(src, file.Name())

		fDestPath := filepath.Join(dest, file.Name())
		if file.IsDir() && recursive {
			os.MkdirAll(fDestPath, 0700)
			Copy(fSrcPath, fDestPath, true, filters)
		} else {
			if filters != nil && len(filters) > 0 {
				isMatch := false

				for _, filter := range filters {
					if strings.Contains(file.Name(), filter) {
						isMatch = true
						break
					}
				}

				if !isMatch {
					continue
				}
			}

			fSrc, err := os.Open(fSrcPath)
			if err != nil {
				return err
			}
			defer fSrc.Close()

			fDest, err := os.OpenFile(
				fDestPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				return err
			}
			defer fDest.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile .
func CopyFile(srcPath, dest string) error {
	fSrc, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fSrc.Close()

	destPath := filepath.Join(dest, filepath.Base(srcPath))
	fDest, err := os.OpenFile(
		destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer fDest.Close()

	_, err = io.Copy(fDest, fSrc)
	if err != nil {
		return err
	}

	return nil
}

// Replace uses Regexp to find any matched from `input` with `regexpTerm`
// and replaces them with `replaceTerm` then returns new string.
func Replace(input *string, regexpTerm string, replaceTerm string) {
	re := regexp.MustCompile(regexpTerm)
	*input = re.ReplaceAllString(*input, replaceTerm)
}

// ModifyFile opens file, changes file content by executing
// `repl` callback function and writes new content.
func ModifyFile(path string, repl func(string) string) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print(err)
		return
	}

	content := repl(string(raw))

	ioutil.WriteFile(path, []byte(content), 0700)
}

// GetSpotifyVersion .
func GetSpotifyVersion(prefsPath string) string {
	pref, err := ini.Load(prefsPath)
	if err != nil {
		log.Fatal(err)
	}

	rootSection, err := pref.GetSection("")
	if err != nil {
		log.Fatal(err)
	}

	version := rootSection.Key("app.last-launched-version")
	return version.MustString("")
}

// GetExecutableDir returns directory of current process
func GetExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Dir(exe)
}

// GetJsHelperDir returns jsHelper directory in executable directory
func GetJsHelperDir() string {
	return filepath.Join(GetExecutableDir(), "jsHelper")
}

// PrependTime prepends current time string to text and returns new string
func PrependTime(text string) string {
	date := time.Now()
	return fmt.Sprintf("%02d:%02d:%02d ", date.Hour(), date.Minute(), date.Second()) + text
}
