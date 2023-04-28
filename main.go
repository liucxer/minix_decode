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
	//__u16 s_imap_blocks;   /* number of inode bitmap blocks */  // 0x0003  inode位图占用3个block       0x1400 - 0x0800 = 0xC00 = 3 * 1024 21 * 1024 inode
	//__u16 s_zmap_blocks;   /* number of zone bitmap blocks */   // 0x0008  data block位图占用8个block  0x3400 - 0x1400 = 0x2000 = 8 * 1024
	//__u16 s_firstdatazone; /* first data zone */                // 0x02b8  = 696
	//__u16 s_log_zone_size; // 0x0000                            // 一个block占用2的0次方K，代表1k
	//__u32 s_max_size;      /* maximum file size */              // 0x10081c00 = 268966912
	//__u16 s_magic;         // 0x138f 魔幻数字
	//__u16 s_state;         // 0x0001          是否挂在，0: 已经挂载， 1: 未挂载
	//__u32 s_zones;         //

	InodeNum             uint16 /* total number of inodes */        // 0x5560  = 21856 inodes
	ZoneNum              uint16 /* total number of zones */         // 0xffff  = 65535 blocks
	InodeBitmapBlocksNum uint16 /* number of inode bitmap blocks */ // 0x0003  inode位图占用3个block       0x1400 - 0x0800 = 0xC00 = 3 * 1024
	ZoneBitmapBlocks     uint16 /* number of zone bitmap blocks */  // 0x0008  data block位图占用8个block  0x3400 - 0x1400 = 0x2000 = 8 * 1024
	FirstDataZone        uint16 /* first data zone */               // 0x02b8  = 696
	ZoneSize             uint16 // 0x0000                            // 一个block占用2的0次方K，代表1k
	MaxFileSize          uint32 /* maximum file size */ // 0x10081c00 = 268966912
	Magic                uint16 // 0x138f 魔幻数字
	State                uint16 // 是否挂在，0: 已经挂载， 1: 未挂载
	Zones                uint32
}

// InodeBitMap  3K大小
type InodeBitMap struct {
	Inode    []uint8
	InodeMap []bool
}

// DataBlockBitMap  8K大小
type DataBlockBitMap struct {
	Inode    []uint8
	InodeMap []bool
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

func DecodeSuperBlock(res *DiskData, bts []byte) error {
	buf := bytes.NewBuffer(bts[1024:2048])
	err := binary.Read(buf, binary.LittleEndian, &res.SuperBlock)
	if err != nil {
		logrus.Errorf("binary.Read err:%v", err)
		return err
	}

	return err
}

func DecodeInodeBitMap(res *DiskData, bts []byte) error {
	buf := bytes.NewBuffer(bts[2048 : 2048+res.SuperBlock.InodeBitmapBlocksNum*1024])
	res.InodeBitMap.Inode = make([]uint8, res.SuperBlock.InodeBitmapBlocksNum*1024)
	err := binary.Read(buf, binary.LittleEndian, &res.InodeBitMap.Inode)
	if err != nil {
		logrus.Errorf("binary.Read err:%v", err)
		return err
	}
	for _, item := range res.InodeBitMap.Inode {
		tmp := make([]bool, 8)
		if item&0x01 == 0x01 {
			tmp[0] = true
		} else {
			tmp[0] = false
		}

		if item&0x02 == 0x02 { // 0000 0010
			tmp[1] = true
		} else {
			tmp[1] = false
		}

		if item&0x04 == 0x04 { // 0000 0100
			tmp[2] = true
		} else {
			tmp[2] = false
		}

		if item&0x08 == 0x08 { // 0000 1000
			tmp[3] = true
		} else {
			tmp[3] = false
		}
		if item&0x10 == 0x10 { // 0001 0000
			tmp[4] = true
		} else {
			tmp[4] = false
		}
		if item&0x20 == 0x20 { // 0010 0000
			tmp[5] = true
		} else {
			tmp[5] = false
		}
		if item&0x40 == 0x40 { // 0100 0000
			tmp[6] = true
		} else {
			tmp[6] = false
		}
		if item&0x80 == 0x80 { // 1000 0000
			tmp[7] = true
		} else {
			tmp[7] = false
		}
		// 0000 0111
		res.InodeBitMap.InodeMap = append(res.InodeBitMap.InodeMap, tmp...)
	}

	return err
}

func DecodeDataBlockBitMap(res *DiskData, bts []byte) error {
	inodeBitMapOffset := 2048 + res.SuperBlock.InodeBitmapBlocksNum*1024
	buf := bytes.NewBuffer(bts[inodeBitMapOffset : inodeBitMapOffset+res.SuperBlock.ZoneBitmapBlocks*1024])
	res.DataBlockBitMap.Inode = make([]uint8, res.SuperBlock.ZoneBitmapBlocks*1024)
	err := binary.Read(buf, binary.LittleEndian, &res.DataBlockBitMap.Inode)
	if err != nil {
		logrus.Errorf("binary.Read err:%v", err)
		return err
	}
	for _, item := range res.DataBlockBitMap.Inode {
		tmp := make([]bool, 8)
		if item&0x01 == 0x01 {
			tmp[0] = true
		} else {
			tmp[0] = false
		}

		if item&0x02 == 0x02 { // 0000 0010
			tmp[1] = true
		} else {
			tmp[1] = false
		}

		if item&0x04 == 0x04 { // 0000 0100
			tmp[2] = true
		} else {
			tmp[2] = false
		}

		if item&0x08 == 0x08 { // 0000 1000
			tmp[3] = true
		} else {
			tmp[3] = false
		}
		if item&0x10 == 0x10 { // 0001 0000
			tmp[4] = true
		} else {
			tmp[4] = false
		}
		if item&0x20 == 0x20 { // 0010 0000
			tmp[5] = true
		} else {
			tmp[5] = false
		}
		if item&0x40 == 0x40 { // 0100 0000
			tmp[6] = true
		} else {
			tmp[6] = false
		}
		if item&0x80 == 0x80 { // 1000 0000
			tmp[7] = true
		} else {
			tmp[7] = false
		}
		// 0000 0111
		res.DataBlockBitMap.InodeMap = append(res.DataBlockBitMap.InodeMap, tmp...)
	}

	return err
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

	// 解析超级块
	err = DecodeSuperBlock(&res, bts)
	if err != nil {
		logrus.Errorf("DecodeSuperBlock err:%v,", err)
		return DiskData{}, err
	}

	// 解析 inode bitmap
	err = DecodeInodeBitMap(&res, bts)
	if err != nil {
		logrus.Errorf("DecodeInodeBitMap err:%v,", err)
		return DiskData{}, err
	}

	// 解析dataBlockBitMap
	err = DecodeDataBlockBitMap(&res, bts)
	if err != nil {
		logrus.Errorf("DecodeDataBlockBitMap err:%v,", err)
		return DiskData{}, err
	}

	return res, err
}

func main() {

	dataPath := "/Users/liuchangxi/Documents/gopath/src/github.com/liucxer/minix_decode/data/data5"
	res, err := DecodeData(dataPath)
	if err != nil {
		logrus.Errorf("DecodeData err:%v", err)
		return
	}
	spew.Dump(res)

}
