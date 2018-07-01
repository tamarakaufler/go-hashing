# go-hashing
Decryption of a file encrypted with unknown hash function and containing a hashed empty line at the beginning. Words are separated by a hash of an empty string. Two formats are considered:
 - lines with character per line
 - files with lines containing incremetally revealed word 

## Implementing a performant solution of finding content of a file with incrementally hashed words

### Problem
Find content of a file with content encrypted with an unknown hash function. The first line of the file is an encrypted empty string. Words are also separated with a line containing an encrypted empty string. Example of a file content with the original content:

 - empty line
 - H
 - He
 - Hel
 - Hell
 - Hello
 - empty line
 - w
 - wo
 - wor
 - worl
 - world
 - empty line (optional)

### Implementation
Concurrency is used to speed up processing. Lines constituting a word are processed concurrently, as is decryption of individual lines. The alphabet and additional symbols are split into several string slices, which are then processed in parallel, trying to find the letter/symbol, that, when added to the previously decryped ones, matches the line hash. On finding a match, the rest of the line processing goroutines is cancelled.

### Usage
go run main.go

go run main.go -f test3.txt
