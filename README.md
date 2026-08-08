# go-logparse

编解码这种小事，犯不着每次都跑在线工具网站溜一圈。

日志分析器，零依赖。能解析常见 combined 格式和「级别+消息」风格日志。

## 用法

```bash
go run . summary access.log        # 按级别计数
go run . status access.log         # 按 HTTP 状态码计数
go run . hourly access.log         # 按小时统计请求数
go run . topips 10 access.log      # 出现最多的 10 个 IP
go run . errors access.log         # 只打印 ERROR/FATAL 行
cat access.log | go run . summary -   # 从管道读
```

级别识别大小写不敏感，缩写 ERR/WARNING 会归一成 ERROR/WARN。时间解析支持 combined 的 `02/Jan/2006:15:04:05 -0700` 和 RFC3339。解析不出来的字段留空，不会报错中断——脏日志也能尽量跑完。
