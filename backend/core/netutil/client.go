package netutil

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	searcherOnce sync.Once
	searcher     *xdb.Searcher
)

// ClientIP extracts the real client IP (X-Forwarded-For > X-Real-IP > RemoteAddr).
func ClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.Index(fwd, ","); i >= 0 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}
	if real := r.Header.Get("X-Real-IP"); real != "" {
		return strings.TrimSpace(real)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// ClientIPFromGin is a gin convenience wrapper.
func ClientIPFromGin(c *gin.Context) string { return ClientIP(c.Request) }

// IsPrivateIP reports whether ip is a private/loopback/link-local address.
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
		return true
	}
	return false
}

// Location resolves ip to a human-readable location via the offline ip2region
// database (config/ip2region.xdb, overridable via IP2REGION_XDB).
// Private IPs -> "内网", unresolvable public IPs -> "外网", so the API never fails.
func Location(ipStr string) string {
	ipStr = strings.TrimSpace(ipStr)
	if ipStr == "" {
		return "-"
	}
	if ipStr == "127.0.0.1" || ipStr == "::1" {
		return "内网"
	}
	if IsPrivateIP(ipStr) {
		return "内网"
	}
	if s := loadSearcher(); s != nil {
		if region, err := s.Search(ipStr); err == nil && region != "" {
			if loc := formatRegion(region); loc != "" {
				return loc
			}
		}
	}
	return "外网"
}

// formatRegion turns "中国|0|广东省|深圳市|电信" into "中国 广东省 深圳市".
func formatRegion(region string) string {
	parts := strings.Split(region, "|")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" || p == "0" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, " ")
}

func loadSearcher() *xdb.Searcher {
	searcherOnce.Do(func() {
		paths := []string{}
		if env := os.Getenv("IP2REGION_XDB"); env != "" {
			paths = append(paths, env)
		}
		paths = append(paths,
			"config/ip2region.xdb",
			"./config/ip2region.xdb",
			"backend/config/ip2region.xdb",
		)
		for _, p := range paths {
			if _, err := os.Stat(p); err != nil {
				continue
			}
			if s, err := xdb.NewWithFileOnly(xdb.IPv4, p); err == nil {
				searcher = s
				return
			}
		}
	})
	return searcher
}

// ParseUA extracts browser/os/device from a User-Agent string (zero-dep heuristic).
func ParseUA(ua string) (browser, osName, device string) {
	l := strings.ToLower(ua)
	// OS
	switch {
	case strings.Contains(l, "windows"):
		osName = "Windows"
	case strings.Contains(l, "android"):
		osName = "Android"
	case strings.Contains(l, "iphone") || strings.Contains(l, "ipad"):
		osName = "iOS"
	case strings.Contains(l, "mac os") || strings.Contains(l, "macintosh"):
		osName = "macOS"
	case strings.Contains(l, "linux"):
		osName = "Linux"
	default:
		osName = "Other"
	}
	// Browser (order matters: Edge UA contains Chrome/Safari too)
	switch {
	case strings.Contains(l, "edg/") || strings.Contains(l, "edge/"):
		browser = "Edge"
	case strings.Contains(l, "opr/") || strings.Contains(l, "opera"):
		browser = "Opera"
	case strings.Contains(l, "firefox/"):
		browser = "Firefox"
	case strings.Contains(l, "chrome/"):
		browser = "Chrome"
	case strings.Contains(l, "safari/") && strings.Contains(l, "version/"):
		browser = "Safari"
	case strings.Contains(l, "msie") || strings.Contains(l, "trident/"):
		browser = "IE"
	default:
		browser = "Other"
	}
	// Device
	switch {
	case strings.Contains(l, "mobile") || strings.Contains(l, "android") ||
		strings.Contains(l, "iphone") || strings.Contains(l, "ipad"):
		device = "Mobile"
	default:
		device = "Desktop"
	}
	return browser, osName, device
}
