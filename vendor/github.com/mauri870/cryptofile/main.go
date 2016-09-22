package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mauri870/cryptofile/crypto"
)

var ()

func main() {
	randomKey := flag.String("key", "", "The key for encryption/decryption")
	in := flag.String("in", "test.txt", "The input file")
	decrypt := flag.Bool("decrypt", false, "Decrypt the file")
	del := flag.Bool("delete", false, "Delete the input file")
	flag.Parse()

	key := checkKeyLength(*randomKey)

	text, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if *del {
			os.Remove(*in)
		}
	}()

	if *decrypt {
		result, err := crypto.Decrypt(key, text)
		if err != nil {
			log.Fatal(err)
		}
		filename := *in
		ioutil.WriteFile(filename[0:len(filename)-len(filepath.Ext(filename))], result, 0600)
		return
	}

	ciphertext, err := crypto.Encrypt(key, text)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(*in+".encrypted", ciphertext, 0600)
}

func checkKeyLength(key string) []byte {
	if key == "" || (len(key) != 32 && len(key) != 16 && len(key) != 24) {
		log.Println("Use -h for usage")
		log.Fatalln(crypto.ErrorKeyInvalidLength)
	}
	return []byte(key)
}
