# go-hashing
Implementing a performant solution of a task finding content of a file with incrementally hashed words

## Problem
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




## Usage

go run main.go
go run main.go -f test3.txt
