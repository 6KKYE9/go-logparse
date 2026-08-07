package main

import "testing"

func TestParseLineLevel(t *testing.T) {
	e := parseLine("2026-08-08 INFO 服务启动")
	if e.Level != "INFO" {
		t.Errorf("级别应识别为 INFO, 得到 %q", e.Level)
	}
	e2 := parseLine("ERROR: 数据库连接失败")
	if e2.Level != "ERROR" {
		t.Errorf("ERROR 应识别, 得到 %q", e2.Level)
	}
	e3 := parseLine("[WARN] 磁盘快满了")
	if e3.Level != "WARN" {
		t.Errorf("WARN 应识别, 得到 %q", e3.Level)
	}
}

func TestParseCommonLog(t *testing.T) {
	line := `127.0.0.1 - - [08/Aug/2026:12:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234`
	e := parseLine(line)
	if e.IP != "127.0.0.1" {
		t.Errorf("IP 解析不对: %q", e.IP)
	}
	if e.Method != "GET" || e.Path != "/api/users" {
		t.Errorf("方法/路径不对: %q %q", e.Method, e.Path)
	}
	if e.Status != 200 {
		t.Errorf("状态码应为 200, 得到 %d", e.Status)
	}
	if e.Bytes != 1234 {
		t.Errorf("字节数应为 1234, 得到 %d", e.Bytes)
	}
	if e.Time.IsZero() {
		t.Error("时间应解析出来")
	}
}

func TestCountByLevel(t *testing.T) {
	entries := []LogEntry{
		{Level: "INFO"}, {Level: "ERROR"}, {Level: "INFO"}, {Level: ""},
	}
	m := countByLevel(entries)
	if m["INFO"] != 2 || m["ERROR"] != 1 || m["UNKNOWN"] != 1 {
		t.Errorf("级别计数不对: %v", m)
	}
}

func TestCountByStatus(t *testing.T) {
	entries := []LogEntry{{Status: 200}, {Status: 404}, {Status: 200}, {Status: 0}}
	m := countByStatus(entries)
	if m[200] != 2 || m[404] != 1 {
		t.Errorf("状态码计数不对: %v", m)
	}
}

func TestTopIPs(t *testing.T) {
	entries := []LogEntry{
		{IP: "1.1.1.1"}, {IP: "1.1.1.1"}, {IP: "2.2.2.2"}, {IP: "1.1.1.1"},
	}
	top := topIPs(entries, 1)
	if top[0].IP != "1.1.1.1" || top[0].Count != 3 {
		t.Errorf("top1 IP 应为 1.1.1.1 x3, 得到 %v", top)
	}
}
