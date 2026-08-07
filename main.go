package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func readLines(path string) ([]string, error) {
	var lines []string
	if path == "-" || path == "" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, scanner.Err()
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		usage()
		return
	}
	cmd := args[0]
	file := args[1]

	lines, err := readLines(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	entries := make([]LogEntry, 0, len(lines))
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		entries = append(entries, parseLine(l))
	}

	switch cmd {
	case "summary":
		levels := countByLevel(entries)
		keys := make([]string, 0, len(levels))
		for k := range levels {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Printf("共解析 %d 行\n", len(entries))
		for _, k := range keys {
			fmt.Printf("  %-8s %d\n", k, levels[k])
		}
	case "status":
		status := countByStatus(entries)
		keys := make([]int, 0, len(status))
		for k := range status {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			fmt.Printf("  %d  %d 次\n", k, status[k])
		}
	case "hourly":
		hours := countByHour(entries)
		for h := 0; h < 24; h++ {
			if hours[h] > 0 {
				fmt.Printf("  %02d 时  %d\n", h, hours[h])
			}
		}
	case "topips":
		n := 5
		if len(args) > 2 {
			fmt.Sscanf(args[2], "%d", &n)
		}
		for _, ip := range topIPs(entries, n) {
			fmt.Printf("  %-15s %d\n", ip.IP, ip.Count)
		}
	case "errors":
		// 只打印 ERROR/FATAL 级别的行
		for _, e := range entries {
			if e.Level == "ERROR" || e.Level == "FATAL" {
				fmt.Println(e.Raw)
			}
		}
	default:
		fmt.Println("不认识子命令:", cmd)
		usage()
	}
}

func usage() {
	fmt.Println("go-logparse 日志分析，零依赖")
	fmt.Println("用法: go-logparse <子命令> <日志文件|->")
	fmt.Println("  summary     各级别计数")
	fmt.Println("  status      各 HTTP 状态码计数")
	fmt.Println("  hourly      按小时统计请求数")
	fmt.Println("  topips [n]  出现最多的 n 个 IP（默认 5）")
	fmt.Println("  errors      只打印 ERROR/FATAL 行")
	fmt.Println("  <日志文件> 用 - 表示从标准输入读")
}
