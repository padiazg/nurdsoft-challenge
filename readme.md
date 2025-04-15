# Nurdsoft coding challenge

## Requirements
-Works on Go 1.24 or above (might work on earlier versions)

## Run the code
```shell
# clone the repo
git clone https://github.com/padiazg/nurdsoft-challenge.git
cd nurdsoft-challenge

# check dependencies
go mod tidy

# build
go build

## run
./nurdsoft-challenge
```

## Call the API
```shell
# add a book
$ curl --header "Content-Type: application/json" --request POST --data '{"Title":"book1", "Author":"author1"}' http://localhost:8080/books
{"ID":1}

# get all books
$ curl --header "Content-Type: application/json" --request GET  http://localhost:8080/books
[{"ID":1,"Title":"book1","Author":"author1","Price":0,"ISBN":"","Active":true},{"ID":2,"Title":"book2","Author":"author2","Price":0,"ISBN":"","Active":true}]

# get a single book by id
$ curl --header "Content-Type: application/json" --request GET  http://localhost:8080/books/1
{"ID":1,"Title":"book1","Author":"author1","Price":0,"ISBN":"","Active":true}

# update a book
$ curl --header "Content-Type: application/json" --request PUT --data '{"Title":"book updated", "Author":"author1"}' http://localhost:8080/books/1
{"ID":1,"Title":"book updated","Author":"author1","Price":0,"ISBN":"","Active":true}

# delete a book
curl --header "Content-Type: application/json" --request DELETE  http://localhost:8080/books/1 
{"delete":"succesfully deleted record 1"}%
```
