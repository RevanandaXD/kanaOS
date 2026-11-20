//go:build !linux

package system

func readDisk(path string) (Disk, bool) {
	return Disk{}, false
}
