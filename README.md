## Get started

通过网络启动，在虚拟机环境下快速安装 Debian 12 。

## Build and Run

执行构建。

```bash
# required go version go1.21.0 and make
make build
```

修改配置后可以直接运行。

```yaml
# 修改成实际机器的网卡
iface: enp0s3
# 设置实际的 IP 地址
ipaddr: 10.0.2.5
logger:
  level: debug
  file: /dev/stdout
dhcp:
  port: 67
tftp:
  port: 69
pxe:
  # 默认打印信息，不安装任何镜像
  default: 0
  entries:
    - label: 0
      display: help
      config: pxelinux.cfg/default
    - label: 1
      display: Debian-12-bookworm-autoinstall
      kernel: images/debian-bookworm-amd64/linux
      initrd: images/debian-bookworm-amd64/initrd.gz
      # 这里的 tftp server 和 ipaddr 的值保持一致，也可以使用外部自定义的 preseed 文件
      append: vga=normal fb=false auto=true priority=critical preseed/url=tftp://10.0.2.5/images/debian-bookworm-amd64/preseed.cfg
```

### Orcale VM VirtualBox

> [!IMPORTANT]
> 1. 网络启动后会停留在镜像选择步骤，选择对应镜像后会在安装过程对硬盘进行格式化，对于已有系统的虚拟电脑应该慎重操作。
> 2. 当前的 Debian 12 安装镜像默认设置普通用户 debian 和 root 用户两个，密码和各自的用户名相同，系统安装完成后请务必修改。

新建虚拟电脑 -> 设置 -> 系统 -> 启动顺序 -> 设置为网络优先。

根据当前服务的运行位置决定网络连接方式：

- 如果已经有多台虚拟电脑运行在同个 NAT 网络下，可以直接将网络设置成 NAT 网络，并在其中任意一台虚拟电脑运行本服务。
- 如果虚拟电脑是独立且全新的，那么直接选择桥接网络，后续将本服务运行在宿主机即可。

安装完成后需要重新调整启动顺序为硬盘优先。

## Todo

- [x] ~~当前的 pxelinux.cfg/default 文件需要渲染 tftp 地址，否则应该手动修改后重新编译。~~
- [x] ~~允许镜像从远程拉取~~。
- [ ] 指定额外的本地 tftp 目录。

## Where file from

### PXE boot

以下所有文件都基于已经安装完成的 Debian 12 系统。

BIOS:

* `pxelinux.0`
* `ldlinux.c32`

``` shell
apt install syslinux
apt install pxelinux
cp /usr/lib/syslinux/modules/bios/ldlinux.c32 tftpboot/ldlinux.c32
cp /usr/lib/PXELINUX/pxelinux.0 tftpboot/pxelinux.0
```

UEFI:

* `syslinux.efi`
* `ldlinux.e64`

``` shell
apt install syslinux
apt install syslinux-efi
cp /lib/syslinux/modules/efi64/ldlinux.e64 tftpboot/ldlinux.e64
cp /usr/lib/SYSLINUX.EFI/efi64/syslinux.efi tftpboot/syslinux.efi
```

### Images

debian-bookworm-amd64:

* `linux`
* `initrd.gz`

``` shell
wget http://http.us.debian.org/debian/dists/bookworm/main/installer-amd64/current/images/netboot/netboot.tar.gz
tar -xvf netboot.tar.gz -C netboot
cp netboot/debian-installer/amd64/linux tftpboot/images/debian-bookworm-amd64/linux
cp netboot/debian-installer/amd64/initrd.gz tftpboot/images/debian-bookworm-amd64/initrd.gz
```

`preseed.cfg` is modified from `example-preseed.txt`.

### Others

* `example-preseed.txt`: download from [d-i.debian.org](https://d-i.debian.org/manual/example-preseed.txt)
