package main

import (
	"bytes"
	"encoding/binary"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// BootBlock 启动块 1K大小
type BootBlock struct {
}

// SuperBlock 超级块 1K大小
type SuperBlock struct {

	//__u16 s_ninodes;       /* total number of inodes */         // 0x5560  = 21856 inodes
	//__u16 s_nzones;        /* total number of zones */          // 0xffff  = 65535 blocks
	//__u16 s_imap_blocks;   /* number of inode bitmap blocks */  // 0x0003  inode位图占用3个block       0x1400 - 0x0800 = 0xC00 = 3 * 1024
	//__u16 s_zmap_blocks;   /* number of zone bitmap blocks */   // 0x0008  data block位图占用8个block  0x3400 - 0x1400 = 0x2000 = 8 * 1024
	//__u16 s_firstdatazone; /* first data zone */                // 0x02b8  = 696
	//__u16 s_log_zone_size; // 0x0000                            // 一个block占用2的0次方K，代表1k
	//__u32 s_max_size;      /* maximum file size */              // 0x10081c00 = 268966912
	//__u16 s_magic;         // 0x138f 魔幻数字
	//__u16 s_state;         // 0x0001          是否挂在，0: 已经挂载， 1: 未挂载
	//__u32 s_zones;         //

	InodeNum             uint16
	ZoneNum              uint16
	InodeBitmapBlocksNum uint16
	ZoneBitmapBlocks     uint16
	FirstDataZone        uint16
	ZoneSize             uint16
	MaxFileSize          uint32
	Magic                uint16
	State                uint16
	Zones                uint32
}

// InodeBitMap  3K大小
type InodeBitMap struct {
}

// DataBlockBitMap  8K大小
type DataBlockBitMap struct {
}

// InodeTable  683K大小
type InodeTable struct {
}

// DataBlock  65535K大小
type DataBlock struct {
}

type DiskData struct {
	BootBlock       BootBlock       `json:"bootBlock"`
	SuperBlock      SuperBlock      `json:"superBlock"`
	InodeBitMap     InodeBitMap     `json:"inodeBitMap"`
	DataBlockBitMap DataBlockBitMap `json:"dataBlockBitMap"`
	InodeTable      InodeTable      `json:"inodeTable"`
	DataBlock       DataBlock       `json:"dataBlock"`
}

func DecodeData(filePath string) (DiskData, error) {
	var (
		err error
		res DiskData
	)

	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("ioutil.ReadFile err:%v, path:%s", err, filePath)
		return DiskData{}, err
	}
	buf := bytes.NewBuffer(bts[1024:2048])
	err = binary.Read(buf, binary.LittleEndian, &res.SuperBlock)
	if err != nil {
		logrus.Errorf("binary.Read err:%v, path:%s", err, filePath)
		return DiskData{}, err
	}
	return res, err
}

func main() {

	dataPath := "/Users/liuchangxi/Documents/gopath/src/github.com/liucxer/minix_decode/data/data2"
	res, err := DecodeData(dataPath)
	if err != nil {
		logrus.Errorf("DecodeData1 err:%v", err)
		return
	}
	spew.Dump(res)

}
