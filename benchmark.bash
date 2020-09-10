#!/bin/bash

cd helpers ; go test -run none -gcflags "-m -m" -bench . -benchtime 3s -benchmem -memprofile ../mem.out 2> ../build.out
