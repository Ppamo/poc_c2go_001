#!/bin/bash

rm -Rf bin
mkdir bin

printf "Building client\n"
APP=app
cd app
printf "> Building for Linux:\n"
GOOS=linux go build -o ../bin/$APP .
if [ $? -eq 0 ]; then
	printf "< Success\n"
else
	printf "< Failure\n"
	exit 1
fi
printf "> Building for Windows:\n"
GOOS=windows go build -o ../bin/$APP.exe .
if [ $? -eq 0 ]; then
	printf "< Success\n"
else
	printf "< Failure\n"
	exit 1
fi
printf "> Building for Mac:\n"
GOOS=darwin go build -o ../bin/$APP.dwn .
if [ $? -eq 0 ]; then
	printf "< Success\n"
else
	printf "< Failure\n"
	exit 1
fi



printf "Building webserver\n"
APP=webserver
cd ../webserver
printf "> Building for Linux:\n"
go build -o ../bin/$APP .
if [ $? -eq 0 ]; then
	printf "< Success\n"
else
	printf "< Failure\n"
	exit 1
fi
printf "> Building for Windows:\n"
GOOS=windows go build -o ../bin/$APP.exe .
if [ $? -eq 0 ]; then
	printf "< Success\n"
else
	printf "< Failure\n"
	exit 1
fi
printf "> Building for Mac:\n"
GOOS=darwin go build -o ../bin/$APP.dwn .
if [ $? -eq 0 ]; then
	printf "< Success\n"
else
	printf "< Failure\n"
	exit 1
fi
