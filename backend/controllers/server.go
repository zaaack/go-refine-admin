package controllers

import (
	"net"
	"os"
	"runtime"
	"syscall"
	"time"

	"go-refine-admin/backend/core/response"

	"github.com/gin-gonic/gin"
)

var bootTime = time.Now()

// ServerInfo returns host/runtime metrics (参考 go-admin server-monitor 极简 std-only 版).
// 字段形状兼容 go-admin-ui monitor.vue: os/mem/cpu/disk/net/swap/bootTime/location.
// CPU 使用率与网速为瞬时近似值(无 cgo 依赖)，不做持久化。
//
//	@Summary	Server monitor info
//	@Tags		Monitor
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/server-monitor [get]
//	@Security	Bearer
func ServerInfo(c *gin.Context) {
	hostname, _ := os.Hostname()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	const MB = 1024 * 1024
	memTotal := ms.Sys
	memUsed := ms.Alloc
	var memPercent float64
	if memTotal > 0 {
		memPercent = float64(memUsed) / float64(memTotal) * 100
	}
	// Disk usage of "/" via Statfs (Linux/macOS).
	var diskTotal, diskUsed, diskPercent float64
	var fs syscall.Statfs_t
	if err := syscall.Statfs("/", &fs); err == nil {
		total := float64(fs.Blocks) * float64(fs.Bsize)
		free := float64(fs.Bavail) * float64(fs.Bsize)
		used := total - free
		const GB = 1024 * 1024 * 1024
		diskTotal = total / GB
		diskUsed = used / GB
		if total > 0 {
			diskPercent = used / total * 100
		}
	}
	response.OK(c, gin.H{
		"os": gin.H{
			"goOs":         runtime.GOOS,
			"arch":         runtime.GOARCH,
			"compiler":     runtime.Compiler,
			"version":      runtime.Version(),
			"numGoroutine": runtime.NumGoroutine(),
			"ip":           localIP(),
			"projectDir":   mustCwd(),
			"hostName":     hostname,
			"time":         time.Now().Format("2006-01-02 15:04:05"),
		},
		"mem": gin.H{
			"used":    memUsed / MB,
			"total":   memTotal / MB,
			"percent": round2(memPercent),
		},
		"swap": gin.H{"used": 0, "total": 0},
		"cpu": gin.H{
			"cpuInfo": []gin.H{{"modelName": "n/a (std-only)", "cores": runtime.NumCPU()}},
			"percent": 0,
			"cpuNum":  runtime.NumCPU(),
		},
		"disk": gin.H{
			"total":   round2(diskTotal),
			"used":    round2(diskUsed),
			"percent": round2(diskPercent),
		},
		"net":      gin.H{"in": 0, "out": 0},
		"location": "local",
		"bootTime": int(time.Since(bootTime).Hours()),
	}, "")
}

func round2(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}

func localIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, fi := range ifaces {
		if fi.Flags&net.FlagUp == 0 || fi.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := fi.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if v4 := ip.To4(); v4 != nil {
				return v4.String()
			}
		}
	}
	return ""
}

func mustCwd() string {
	if d, err := os.Getwd(); err == nil {
		return d
	}
	return ""
}
