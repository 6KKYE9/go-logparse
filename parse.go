package main

import (
	"regexp"
	"strings"
	"time"
)

// 常见日志行的正则，尽量兼容 nginx/apache 风格和简单 "时间 [级别] 消息"。
var commonLog = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+\[([^\]]+)\]\s+"(\S+)\s+(\S+)\s+(\S+)"\s+(\d+)\s+(\d+|-)`)

// LogEntry 是一条解析出来的日志。
type LogEntry struct {
	Raw     string
	Level   string
	Time    time.Time
	Message string
	IP      string
	Method  string
	Path    string
	Status  int
	Bytes   int
}

// 尝试从一行里抠出级别和时间，抠不到就留空。
func parseLine(line string) LogEntry {
	e := LogEntry{Raw: line}

	// 级别：像 [ERROR level] 或 "level=info" 或行首的 ERROR/WARN/INFO/DEBUG
	if m := regexp.MustCompile(`(?i)\b(ERROR|ERR|WARN|WARNING|INFO|DEBUG|TRACE|FATAL)\b`).FindString(line); m != "" {
		e.Level = strings.ToUpper(m)
	}
	// 级别缩写映射
	switch e.Level {
	case "ERR":
		e.Level = "ERROR"
	case "WARNING":
		e.Level = "WARN"
	}

	// 时间：先找 [xxx] 里的，再找 RFC3339 或 常见格式
	if m := regexp.MustCompile(`\[([^\]]+)\]`).FindStringSubmatch(line); m != nil {
		if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", m[1]); err == nil {
			e.Time = t
		} else if t, err := time.Parse(time.RFC3339, m[1]); err == nil {
			e.Time = t
		}
	}

	// 常见 combined 格式
	if m := commonLog.FindStringSubmatch(line); m != nil {
		e.IP = m[1]
		e.Method = m[5]
		e.Path = m[6]
		e.Status = atoiSafe(m[8])
		if m[9] != "-" {
			e.Bytes = atoiSafe(m[9])
		}
	}

	// 消息：去掉前面解析过的部分，剩下的当消息
	e.Message = line
	return e
}

func atoiSafe(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return n
		}
		n = n*10 + int(r-'0')
	}
	return n
}

// 按级别统计数量。
func countByLevel(entries []LogEntry) map[string]int {
	m := map[string]int{}
	for _, e := range entries {
		lvl := e.Level
		if lvl == "" {
			lvl = "UNKNOWN"
		}
		m[lvl]++
	}
	return m
}

// 按状态码统计。
func countByStatus(entries []LogEntry) map[int]int {
	m := map[int]int{}
	for _, e := range entries {
		if e.Status != 0 {
			m[e.Status]++
		}
	}
	return m
}

// 按小时统计请求数（需要时间能解析出来）。
func countByHour(entries []LogEntry) map[int]int {
	m := map[int]int{}
	for _, e := range entries {
		if !e.Time.IsZero() {
			m[e.Time.Hour()]++
		}
	}
	return m
}

// 找出最频繁的 IP。
func topIPs(entries []LogEntry, n int) []struct {
	IP    string
	Count int
} {
	m := map[string]int{}
	for _, e := range entries {
		if e.IP != "" {
			m[e.IP]++
		}
	}
	pairs := make([]ipCount, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, ipCount{k, v})
	}
	sortPairs(pairs)
	if n > len(pairs) {
		n = len(pairs)
	}
	out := make([]struct {
		IP    string
		Count int
	}, n)
	for i := 0; i < n; i++ {
		out[i] = pairs[i]
	}
	return out
}

// ipCount 是 IP 和出现次数的配对，topIPs 和 sortPairs 共用。
type ipCount struct {
	IP    string
	Count int
}

func sortPairs(pairs []ipCount) {
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].Count > pairs[i].Count {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
}
