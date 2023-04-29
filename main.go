package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/minix_decode/decoder"
	"github.com/sirupsen/logrus"
)

func main() {
	dataPath := "/Users/liuchangxi/Documents/gopath/src/github.com/liucxer/minix_decode/data/data7.dat"
	var disk decoder.DiskData
	err := disk.Decode(dataPath)
	if err != nil {
		logrus.Errorf("disk.Decode err:%v", err)
		return
	}

	spew.Dump(disk)
}
