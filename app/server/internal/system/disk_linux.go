//go:build linux

package system

import "syscall"

func readDisk(path string) (Disk, bool) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return Disk{}, false
	}
	bsize := uint64(st.Bsize)
	total := uint64(st.Blocks) * bsize
	free := uint64(st.Bavail) * bsize
	if total == 0 || free > total {
		return Disk{}, false
	}
	return Disk{Path: path, Total: total, Used: total - free}, true
}
