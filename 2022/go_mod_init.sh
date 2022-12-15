#!/usr/bin/env bash

go mod init github.com/davidjcmar/adventofcode/$(basename $PWD)
touch main.go input.txt test_input.txt
