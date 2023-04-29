package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/minix_decode/decoder"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	var disk decoder.DiskData
	err := disk.Decode(os.Args[1])
	if err != nil {
		logrus.Errorf("disk.Decode err:%v", err)
		return
	}

	spew.Dump(disk)
}
