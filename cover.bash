#!/bin/bash

cd helpers ; go test . -coverprofile ../cover.out -trace ../trace.out -v

go tool cover -html=../cover.out
