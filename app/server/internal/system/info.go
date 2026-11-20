package system

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type Disk struct {
	Path  string `json:"path"`
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

type Info struct {
	OS       string `json:"os"`
	Hostname string `json:"hostname"`
	Kernel   string `json:"kernel"`
	Arch     string `json:"arch"`
	Uptime   uint64 `json:"uptime"`
	CPUs     int    `json:"cpus"`
	CPUModel string `json:"cpu_model"`
	MemTotal uint64 `json:"mem_total"`
	MemUsed  uint64 `json:"mem_used"`
	Disks    []Disk `json:"disks"`
}

func Collect() Info {
	host, _ := os.Hostname()
	info := Info{
		OS:       runtime.GOOS,
		Hostname: host,
		Arch:     runtime.GOARCH,
		CPUs:     runtime.NumCPU(),
	}

	if name := readOSRelease(); name != "" {
		info.OS = name
	}
	if kernel, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		info.Kernel = strings.TrimSpace(string(kernel))
	}
	if uptime := readUptime(); uptime > 0 {
		info.Uptime = uptime
	}
	if model := readCPUModel(); model != "" {
		info.CPUModel = model
	}
	if total, used, ok := readMemInfo(); ok {
		info.MemTotal = total
		info.MemUsed = used
	}
	if disk, ok := readDisk("/"); ok {
		info.Disks = []Disk{disk}
	}

	return info
}

func readOSRelease() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return ""
	}
	defer f.Close()

	var pretty, name string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if v, ok := cutPrefix(line, "PRETTY_NAME="); ok {
			pretty = unquote(v)
		} else if v, ok := cutPrefix(line, "NAME="); ok && name == "" {
			name = unquote(v)
		}
	}
	if pretty != "" {
		return pretty
	}
	return name
}

func readUptime() uint64 {
	b, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(b))
	if len(fields) == 0 {
		return 0
	}
	sec, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0
	}
	return uint64(sec)
}

func readCPUModel() string {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		for _, key := range []string{"model name", "Hardware", "Model", "Processor"} {
			if strings.HasPrefix(line, key) {
				if _, v, ok := strings.Cut(line, ":"); ok {
					v = strings.TrimSpace(v)
					if v != "" {
						return v
					}
				}
			}
		}
	}
	return ""
}

func readMemInfo() (total, used uint64, ok bool) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, false
	}
	defer f.Close()

	var memTotal, memAvailable uint64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if v, ok := memValue(line, "MemTotal:"); ok {
			memTotal = v
		} else if v, ok := memValue(line, "MemAvailable:"); ok {
			memAvailable = v
		}
		if memTotal > 0 && memAvailable > 0 {
			break
		}
	}
	if memTotal == 0 {
		return 0, 0, false
	}
	if memAvailable > memTotal {
		memAvailable = memTotal
	}
	return memTotal, memTotal - memAvailable, true
}

func memValue(line, key string) (uint64, bool) {
	if !strings.HasPrefix(line, key) {
		return 0, false
	}
	fields := strings.Fields(strings.TrimPrefix(line, key))
	if len(fields) == 0 {
		return 0, false
	}
	v, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, false
	}
	// /proc/meminfo values are in kB
	return v * 1024, true
}

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

func cutPrefix(s, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}
	return "", false
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
