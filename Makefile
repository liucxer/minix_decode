#----------------------------------------------
#    Name: Makefile
#  Author: liucx
# Version: 1.*
#   Usage: build minix_decoder

#使用说明:
#1.编译
#make
#或 make all
#2.分发
#make install ip=10.0.8.126 pass=daemon
#3.清理
#make clean
#--------------------Declare-------------------

# const
IP := $(ip)
SERVER_PASSWORD := $(pass)
SCP_STORAGE_PATH := /root

all: minix_decoder minix_decoder_mac minix_decoder_linux_x86_64 minix_decoder_linux_aarch64 minix_decoder_windows

minix_decoder:
	mkdir -p build && cd build
	go build -o ./build/minix_decoder

minix_decoder_mac:
	mkdir -p build && cd build
	go build -o ./build/minix_decoder_mac

minix_decoder_linux_x86_64:
	mkdir -p build && cd build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/minix_decoder_linux_x86_64

minix_decoder_linux_aarch64:
	mkdir -p build && cd build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/minix_decoder_linux_aarch64

minix_decoder_windows:
	mkdir -p build && cd build
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/minix_decoder_windows

clean:
	\rm -rf ./build/*

install:
	make clean
	make minix_decoder
	sshpass -p "$(SERVER_PASSWORD)" scp -o StrictHostKeyChecking=no build/minix_decoder root@$(IP):$(SCP_STORAGE_PATH)/