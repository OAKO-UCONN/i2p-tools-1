// +build linux
// tun_linux.go -- tun interface with cgo for linux / bsd
//

package samtun

/*

#include <string.h>
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <arpa/inet.h>
#include <sys/ioctl.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <linux/if.h>
#include <linux/if_tun.h>

int tundev_open(char * ifname) {
  int fd = open("/dev/net/tun", O_RDWR);
  if (fd > 0) {
    struct ifreq ifr;
    memset(&ifr, 0, sizeof(ifr));
    ifr.ifr_flags = IFF_TUN | IFF_NO_PI;
    strncpy(ifr.ifr_name, ifname, IFNAMSIZ);
    if ( ioctl(fd , TUNSETIFF, (void*) &ifr) < 0) {
      close(fd);
      return -1;
    }
  }
  return fd;
}

int tundev_up(char * ifname, char * addr, char * dstaddr, int mtu) {

  struct ifreq ifr;
  memset(&ifr, 0, sizeof(ifr));
  strncpy(ifr.ifr_name, ifname, IFNAMSIZ);
  int fd = socket(AF_INET, SOCK_DGRAM, IPPROTO_IP);
  if ( fd > 0 ) {
    if ( ioctl(fd, SIOCGIFINDEX, (void*) &ifr) < 0 ) {
      perror("SIOCGIFINDEX");
      close(fd);
      return -1;
    }
    int idx = ifr.ifr_ifindex;
    ifr.ifr_mtu = mtu;
    if ( ioctl(fd, SIOCSIFMTU, (void*) &ifr) < 0) {
      close(fd);
      perror("SIOCSIFMTU");
      return -1;
    }

    struct sockaddr_in dst;
    memset(&dst, 0, sizeof(struct sockaddr_in));
    dst.sin_family = AF_INET;
    if ( ! inet_aton(dstaddr, &dst.sin_addr) ) {
      printf("invalid dstaddr %s\n", dstaddr);
      close(fd);
      return -1;
    }
    struct sockaddr_in src;
    memset(&src, 0, sizeof(struct sockaddr_in));
    src.sin_family = AF_INET;
    if ( ! inet_aton(addr, &src.sin_addr) ) {
      printf("invalid srcaddr %s\n", addr);
      close(fd);
      return -1;
    }

    memcpy(&ifr.ifr_addr, &src, sizeof(struct sockaddr_in));
    if ( ioctl(fd, SIOCSIFADDR, (void*)&ifr) < 0 ) {
      close(fd);
      perror("SIOCSIFADDR");
     return -1;
    }
    memcpy(&ifr.ifr_dstaddr, &dst, sizeof(struct sockaddr_in));

    if ( ioctl(fd, SIOCSIFDSTADDR, (void*)&ifr) < 0 ) {
      close(fd);
      perror("SIOCSIFDSTADDR");
     return -1;
    }
    if ( ioctl(fd, SIOCGIFFLAGS, (void*)&ifr) < 0 ) {
      close(fd);
      perror("SIOCGIFFLAGS");
      return -1;
    }
    ifr.ifr_flags |= IFF_UP | IFF_RUNNING;
    if ( ioctl(fd, SIOCSIFFLAGS, (void*)&ifr) < 0 ) {
      close(fd);
      return -1;
    }
    close(fd);
    return 0;
  } 
  return -1;
}

void tundev_close(int fd) {
  close(fd);
}

*/
import "C"

import (
  "errors"
)

type tunDev C.int

func newTun(ifname, addr, dstaddr string, mtu int) (t tunDev, err error) {
  fd := C.tundev_open(C.CString(ifname))
  
  if fd == -1 {
    err = errors.New("cannot open tun interface")
  } else {
    if C.tundev_up(C.CString(ifname), C.CString(addr), C.CString(dstaddr), C.int(mtu)) < C.int(0) {
      err = errors.New("cannot put up interface")
    } else {
      return tunDev(fd), nil
    }
  }
  return -1, err
}

// read from the tun device
func (t tunDev) Read(d []byte) (n int, err error) {
  return fdRead(C.int(t), d)
}

func (t tunDev) Write(d []byte) (n int, err error) {
  return fdWrite(C.int(t), d)
}


func (t tunDev) Close() {
  fdClose(C.int(t))
}
