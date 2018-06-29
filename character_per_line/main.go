package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var filename string

func init() {
	flag.StringVar(&filename, "f", "test.txt", "File name whose content you want to decipher.")
}

func main() {
	flag.Parse()
	fmt.Println("Solve the puzzle - decrypt the file with encrypted with unknown hash function")

	start := time.Now()

	b, err := ioutil.ReadFile(filename)
	checkError(err)
	// process
	s := string(b)
	lines := strings.Split(s, "\n")
	esh := lines[0]
	alg, decryptM, err := createDecryptMap(esh)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for k, v := range decryptM {
		fmt.Printf("%s ... %s\n", k, v)
	}

	emptyB := sha1.Sum([]byte(""))
	fmt.Printf("sha1.Sum... %x\n", string(emptyB[:]))

	var decrypted string
	if decrypted, err = decryptLines(esh, decryptM, lines); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("File hashed with %s algorithm\n", alg)
	fmt.Printf("File content: %s\n", decrypted)

	fmt.Printf("Took %s\n", time.Since(start))
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createDecryptMap(emptyString string) (string, map[string]string, error) {
	decrypt := make(map[string]string)
	alpha := "a b c d e f g h i j k l m n o p q r s t u v w x y z A B C D E F G H I J K L M N O P Q R S T U V W X Y Z 0 1 2 3 4 5 6 7 8 9"
	alphaList := strings.Split(alpha, " ")

	alg := ""
	switch emptyString {
	case "da39a3ee5e6b4b0d3255bfef95601890afd80709":
		alg = "sha1"

		lEncr := sha1.New()
		lEncr.Write([]byte(""))
		digest := lEncr.Sum(nil)
		lEncrHex := fmt.Sprintf("%x", digest)
		fmt.Printf("==> %x", digest)
		decrypt[lEncrHex] = ""

		for _, l := range alphaList {
			lEncr := sha1.Sum([]byte(l))
			//fmt.Printf("%d -> %x\n", len(emptyString), string(lEncr[:]))
			lEncrHex := fmt.Sprintf("%x", string(lEncr[:]))
			decrypt[lEncrHex] = l
		}
	case "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855":
		alg = "sha256"
		for _, l := range alphaList {
			lEncr := sha256.Sum256([]byte(l))
			lEncrHex := fmt.Sprintf("%x", string(lEncr[:]))
			decrypt[lEncrHex] = l
		}
	case "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e":
		alg = "sha512"
		for _, l := range alphaList {
			lEncr := sha512.Sum512([]byte(l))
			lEncrHex := fmt.Sprintf("%x", string(lEncr[:]))
			decrypt[lEncrHex] = l
		}
	case "d41d8cd98f00b204e9800998ecf8427e":
		alg = "md5"
		for _, l := range alphaList {
			lEncr := md5.Sum([]byte(l))
			lEncrHex := fmt.Sprintf("%x", string(lEncr[:]))
			decrypt[lEncrHex] = l
		}
	default:
		return "", nil, errors.New("Unknown algorithm")
	}

	return alg, decrypt, nil
}

func decryptLines(emptyEncr string, decrypt map[string]string, lines []string) (string, error) {
	content := ""
	temp := []string{}
	for i, l := range lines {
		if i == 0 {
			continue
		}
		if l == emptyEncr || i == len(lines)-1 {
			if i == len(lines)-1 {
				temp = append(temp, l)
			}
			for _, t := range temp {
				if t != emptyEncr {
					d, ok := decrypt[t]
					if !ok {
						return "", fmt.Errorf("Incorrect encryption, unknown: %s\n", t)
					}
					content = content + d
				}
			}
			if l == emptyEncr && i != len(lines)-1 {
				content += " "
			}
			temp = []string{}
			continue
		}
		temp = append(temp, l)
	}
	return content, nil
}
