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
	"sync"
	"time"
)

var filename string
var concur int

type result struct {
	pos  int
	word string
}

type encryptFunc func(string) string

func init() {
	flag.StringVar(&filename, "f", "test.txt", "File name whose content you want to decipher. Default is test.txt")
	flag.IntVar(&concur, "n", 4, "Concurrency. Default is 4.")
}

func main() {
	flag.Parse()
	fmt.Println("Solve the puzzle - decrypt the file with encrypted with unknown hash function")

	start := time.Now()

	b, err := ioutil.ReadFile(filename)
	checkError(err)

	// Determine the encryption type and create a helper to find an encrypted letter
	s := string(b)
	lines := strings.Split(s, "\n")
	esh := lines[0]
	alg, decryptM, err := createDecryptMap(esh)
	if err != nil {
		fmt.Println(alg, err)
		os.Exit(1)
	}

	var encryptF encryptFunc
	switch alg {
	case "sha1":
		encryptF = encryptSHA1()
	case "sha256":
		encryptF = encryptSHA256()
	case "sha512":
		encryptF = encryptSHA512()
	case "md5":
		encryptF = encryptMD5()
	default:
		fmt.Printf("Unknown algorithm")
		os.Exit(1)
	}

	lines = lines[1:]
	words := separateWords(concur, esh, lines)

	// Process each word concurrently
	// ==============================
	//  => gather and present the file content
	var wg sync.WaitGroup // wait for all word processing goroutines

	alphaConc := createLetterSlices(concur)
	resultCh := make(chan result, len(words))

	//---------------------------------------------------------------
	for wordPos, letterLines := range words {
		wg.Add(1)
		go func(wordPos int, letterLines []string) {
			defer wg.Done()
			firstLetter := decryptM[letterLines[0]]
			fmt.Printf("first letter = %v\n", firstLetter)

			processWord(
				wordPos, resultCh,
				encryptF, alphaConc, firstLetter, letterLines[1:])

		}(wordPos, letterLines)
	}

	// wait until all words are processed
	wg.Wait()
	close(resultCh)
	//---------------------------------------------------------------

	contentM := map[int]string{}
	content := []string{}

	for res := range resultCh {
		fmt.Printf("\nResult %#v\n\n", res)
		contentM[res.pos] = res.word
	}
	for i := 0; i < len(contentM); i++ {
		content = append(content, contentM[i])
	}
	fmt.Printf("File content: %s\n", strings.Join(content, " "))
	fmt.Printf("Took %s\n", time.Since(start))
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// separateWords puts lines denoting one word into separate slices to allow concurrent processing.
func separateWords(n int, esh string, lines []string) map[int][]string {
	wordLines := map[int][]string{}
	wn := 0
	for _, wl := range lines {
		// separate words on empty line
		if wl == esh {
			wn++
			continue
		}
		wordLines[wn] = append(wordLines[wn], wl)
	}
	return wordLines
}

// createLetterSlices segments the alphabet into several sub slices to allow concurrent
// search.
func createLetterSlices(n int) map[int][]string {
	letterMap := make(map[int][]string)
	alpha := "a b c d e f g h i j k l m n o p q r s t u v w x y z A B C D E F G H I J K L M N O P Q R S T U V W X Y Z"
	alphaList := strings.Split(alpha, " ")

	alphaListLen := len(alphaList)
	concDelta := alphaListLen / n
	conc := concDelta
	concKey := 0
	for i := 0; i < alphaListLen; i++ {
		if i < conc {
			letterMap[concKey] = append(letterMap[concKey], alphaList[i])
			continue
		}
		conc = conc + concDelta
		concKey++
		letterMap[concKey] = append(letterMap[concKey], alphaList[i])
	}

	return letterMap
}

func createDecryptMap(emptyString string) (string, map[string]string, error) {
	decrypt := make(map[string]string)
	alpha := "a b c d e f g h i j k l m n o p q r s t u v w x y z A B C D E F G H I J K L M N O P Q R S T U V W X Y Z 0 1 2 3 4 5 6 7 8 9"
	alphaList := strings.Split(alpha, " ")

	alg := ""
	switch emptyString {
	case "da39a3ee5e6b4b0d3255bfef95601890afd80709":
		alg = "sha1"
		for _, l := range alphaList {
			lEncr := sha1.Sum([]byte(l))
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

func encryptSHA1() encryptFunc {
	return func(plaintext string) string {
		return fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	}
}
func encryptSHA256() encryptFunc {
	return func(plaintext string) string {
		return fmt.Sprintf("%x", sha256.Sum256([]byte(plaintext)))
	}
}
func encryptSHA512() encryptFunc {
	return func(plaintext string) string {
		return fmt.Sprintf("%x", sha512.Sum512([]byte(plaintext)))
	}
}
func encryptMD5() encryptFunc {
	return func(plaintext string) string {
		return fmt.Sprintf("%x", md5.Sum([]byte(plaintext)))
	}
}

// processWord processes a slice of lines incrementally representing one word.
// Eg:
// h
// ha
// hap
// happ
// happy
func processWord(wordPos int, resultCh chan result,
	encryptF encryptFunc, alphaConc map[int][]string, firstLetter string, wordLines []string) {

	decrypted := firstLetter
	decryptedCh := make(chan string, 1)
	concCount := len(alphaConc)
	done := make(chan struct{}, concCount)

	// encrypted is the original encrypted line content
	lastLineIdx := len(wordLines) - 1
	for lineIdx := 0; lineIdx <= lastLineIdx; lineIdx++ {

		fmt.Printf("\n==> processWord: processing line %d\n", lineIdx)

		var wg sync.WaitGroup
		encryptedLine := wordLines[lineIdx]
		fmt.Printf("Decrypting [%s, %d]\n", encryptedLine, lastLineIdx)

		// concurrent processing of one line witnin the slice of lines representing a word
		for i, letters := range alphaConc {
			fmt.Printf("Processing %d [%#v]\n", i, letters)

			wg.Add(1)
			go func(letters []string) {
				defer wg.Done()
				decipherLine(done, decryptedCh, concCount, encryptF, letters, decrypted, encryptedLine)
			}(letters)
		}

		wg.Wait()

		fmt.Println("waiting to receive from decryptedCh representing partial word\n")
		select {
		case decrypted = <-decryptedCh:
			fmt.Printf(">>> RECEIVED decrypted line within the word = %s\n", decrypted)

		case <-time.After(3 * time.Second):
			close(decryptedCh)

			return
		}

		if lineIdx == lastLineIdx {
			close(decryptedCh)
			resultCh <- result{
				pos:  wordPos,
				word: decrypted,
			}
		}
	}

}

// decipherLine decrypts a line.
func decipherLine(done chan struct{}, decryptedCh chan string, concCount int, encryptF encryptFunc, letters []string, decrypted string, encrypted string) {
	for _, l := range letters {
		trial := decrypted + l
		if encryptF(trial) == encrypted {
			fmt.Printf("\t\tIn decipherLine: SENDING [%s] using %s and %v \n", trial, decrypted, letters)
			decryptedCh <- trial
			return
		}
		trial = decrypted
	}
}
