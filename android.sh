#!/bin/bash

fyne-cross android -arch arm64 -app-id cn.poem.flower -ldflags "-s -w" -env GOPROXY=https://goproxy.cn,direct
