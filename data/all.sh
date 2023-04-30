## 完全清空
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
dd if=/dev/vdb of=./data1.dat bs=1M count=1

## 只完成minix格式化
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
mkfs.minix /dev/vdb
dd if=/dev/vdb of=./data2.dat bs=1M count=1

## 完成minix格式化, 创建文件/1
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
mkfs.minix /dev/vdb
mount /dev/vdb /mnt
touch /mnt/1
umount /mnt
dd if=/dev/vdb of=./data3.dat bs=1M count=1

## 完成minix格式化, 创建文件/1, 写入内容11
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
mkfs.minix /dev/vdb
mount /dev/vdb /mnt
touch /mnt/1
echo "11" > /mnt/1
umount /mnt
dd if=/dev/vdb of=./data4.dat bs=1M count=1

# 完成minix格式化,
# 创建文件/1, 写入内容11,
# 创建目录2
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
mkfs.minix /dev/vdb
mount /dev/vdb /mnt
touch /mnt/1
echo "11" > /mnt/1
mkdir -p /mnt/2
umount /mnt
dd if=/dev/vdb of=./data5.dat bs=1M count=1

# 完成minix格式化,
# 创建文件/1, 写入内容11,
# 创建目录2,
# 创建文件/2/3
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
hexdump -C /dev/vdb
mkfs.minix /dev/vdb
mount /dev/vdb /mnt
touch /mnt/1
echo "11" > /mnt/1
mkdir -p /mnt/2
touch /mnt/2/3
umount /mnt
dd if=/dev/vdb of=./data6.dat bs=1M count=1

# 完成minix格式化,
# 创建文件/1, 写入内容11,
# 创建目录2,
# 创建文件/2/3, 写入内容33
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
mkfs.minix /dev/vdb
mount /dev/vdb /mnt

touch /mnt/1
echo "11" > /mnt/1
chmod 777 /mnt/1

mkdir -p /mnt/2
chmod 777 /mnt/2

touch /mnt/2/3
echo "33" > /mnt/2/3
chmod 777 /mnt/2/3

umount /mnt
dd if=/dev/vdb of=./data7.dat bs=1M count=1

