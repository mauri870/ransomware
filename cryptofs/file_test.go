package cryptofs

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "us55VKQ3GT7K1YxHwFyO1KS7ULC14ILo"

	pwd, err := os.Getwd()
	handleErr(err, t)

	err = ioutil.WriteFile("original.txt", []byte("Hello World"), 0600)
	handleErr(err, t)

	f, err := os.Open("original.txt")
	handleErr(err, t)
	defer func() {
		f.Close()
		os.Remove("original.txt")
	}()

	fstat, err := f.Stat()
	handleErr(err, t)

	fileInfo := &File{fstat, "txt", pwd + `/original.txt`}

	var buf []byte
	buffer := bytes.NewBuffer(buf)
	err = fileInfo.Encrypt(key, buffer)
	handleErr(err, t)

	content, err := ioutil.ReadAll(buffer)
	handleErr(err, t)

	err = ioutil.WriteFile("enc.txt", content, 0600)
	handleErr(err, t)

	f2, err := os.Open("enc.txt")
	handleErr(err, t)
	defer func() {
		f2.Close()
		os.Remove("enc.txt")
	}()

	fstat, err = f2.Stat()
	handleErr(err, t)

	fileInfo = &File{fstat, "txt", pwd + `/enc.txt`}

	buffer.Reset()
	err = fileInfo.Decrypt(key, buffer)

	content, err = ioutil.ReadAll(buffer)
	handleErr(err, t)

	if string(content) != "Hello World" {
		t.Errorf("Expect 'Hello World' but got %s", string(content))
	}

}

func handleErr(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("An error ocurred: %s", err)
	}
}
