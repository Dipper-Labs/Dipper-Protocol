package cli

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCodeFromFile(t *testing.T) {
	codeFile, err := filepath.Abs("/Users/a1/go/src/github.com/Dipper-Protocol/demo/storage.bc")
	//if 0 == len(codeFile) {
	//	return nil, errors.New("code_file can not be empty")
	//}
	fmt.Println(codeFile)
	fmt.Println(len(codeFile))
	fmt.Println("===0===")
	hexcode, err := ioutil.ReadFile(codeFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hexcode)
	fmt.Println(len(hexcode))
	fmt.Println("===1==")
	hexcode = bytes.TrimSpace(hexcode)
	fmt.Println(hexcode)
	fmt.Println(len(hexcode))

	fmt.Println("===2====")
	if len(hexcode)%2 != 0 {
		//return nil, errors.New(fmt.Sprintf("Invalid input length for hex data (%d)\n", len(hexcode)))
	}
	fmt.Println(string(hexcode))

	code, err := hex.DecodeString(string(hexcode))
	fmt.Println(code, err)
	if err != nil {
		//return nil, err
	}


	//return code, nil
	//code, err := CodeFromFile("/Users/a1/go/src/github.com/Dipper-Protocol/demo/storage.bc")
	//fmt.Println(code, err)
}