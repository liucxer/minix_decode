package decoder

import (
	"bytes"
	"encoding/binary"
	"github.com/sirupsen/logrus"
)

// Inode  683K大小
type Inode struct {
	Mode   uint16    // 文件类型和属性(rwx 位)。
	Uid    uint16    // 用户id（文件拥有者标识符）。
	Size   uint32    // 文件大小（字节数）。
	Time   uint32    // 修改时间（自1970.1.1:0 算起，秒）。
	Gid    uint8     // 组id(文件拥有者所在的组)。
	NLinks uint8     // 链接数（多少个文件目录项指向该i 节点）。
	Zone   [9]uint16 // 直接(0-6)、间接(7)或双重间接(8)逻辑块号。
}

type InodeItem struct {
	Inode
	Data []byte `json:"data"`
}

func (inode *InodeItem) ReadData(bts []byte) error {
	var err error
	for i := 0; i < 7; i++ {
		if inode.Zone[i] == 0 {
			break
		}
		start := int64(inode.Zone[i]) * 1024
		end := int64(inode.Zone[i])*1024 + 1024
		data := bts[start:end]
		inode.Data = append(inode.Data, data...)
	}

	if inode.Zone[7] != 0 {
		start := inode.Zone[7] * 1024
		end := inode.Zone[7]*1024 + 1024
		inodeTable := bts[start:end]
		zoneList := make([]uint16, 1024/16)
		buf := bytes.NewBuffer(inodeTable)
		err := binary.Read(buf, binary.LittleEndian, &zoneList)
		if err != nil {
			logrus.Errorf("binary.Read err:%v", err)
			return err
		}

		for _, item := range zoneList {
			if item == 0 {
				break
			}
			start := item * 1024
			end := item*1024 + 1024
			inode.Data = append(inode.Data, bts[start:end]...)
		}
	}

	return err
}

type InodeTable struct {
	InodeItems []InodeItem `json:"inodeItems"`
}

func (inodeTable *InodeTable) Decode(bts []byte, allBts []byte, inodeNum int64) error {
	var (
		err error
	)

	buf := bytes.NewBuffer(bts)
	inodeList := make([]Inode, inodeNum)
	err = binary.Read(buf, binary.LittleEndian, inodeList)
	if err != nil {
		logrus.Errorf("binary.Read err:%v", err)
		return err
	}
	for _, inode := range inodeList {
		var inodeItem InodeItem
		inodeItem.Inode = inode
		inodeTable.InodeItems = append(inodeTable.InodeItems, inodeItem)
	}

	for i := 0; i < len(inodeTable.InodeItems); i++ {
		err = (&inodeTable.InodeItems[i]).ReadData(allBts)
		if err != nil {
			logrus.Errorf("inode.ReadData err:%v", err)
			return err
		}
	}

	return err
}
