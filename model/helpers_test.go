package model

import (
	"bou.ke/monkey"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func PatchNoHomeDir() {
	monkey.Patch(os.UserHomeDir, func() (string, error) {
		return "", fmt.Errorf("boom")
	})
}

func PatchHomeDir(dir string) {
	monkey.Patch(os.UserHomeDir, func() (string, error) {
		return dir, nil
	})
}

func PatchJsonMarshalFail() {
	monkey.Patch(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("json marshal fail")
	})
}

func PatchJsonUnmarshalFail() {
	monkey.Patch(json.Unmarshal, func(b []byte, v interface{}) error {
		return fmt.Errorf("json unmarshal fail")
	})
}

func PatchWriteFileFail() {
	monkey.Patch(ioutil.WriteFile, func(filename string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("write file fail")
	})
}

var fileRedirectGuard *monkey.PatchGuard

func PatchWriteFileRedirect(redirectTo string) {
	fileRedirectGuard = monkey.Patch(ioutil.WriteFile, func(filename string, data []byte, perm os.FileMode) error {
		fileRedirectGuard.Unpatch()
		defer fileRedirectGuard.Restore()

		return ioutil.WriteFile(redirectTo, data, perm)
	})
}

func PatchOpenFail() {
	monkey.Patch(os.Open, func(filename string) (*os.File, error) {
		return nil, fmt.Errorf("open fail")
	})
}

func PatchReadAllFail() {
	monkey.Patch(ioutil.ReadAll, func(r io.Reader) ([]byte, error) {
		return nil, fmt.Errorf("read all fail")
	})
}
