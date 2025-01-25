#!/bin/bash

# 設定 pprof 伺服器的地址
pprof_url="http://localhost:6060/debug/pprof"
pprof_file="/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/pprof-Output"

# 生成各種性能報告
curl -o $pprof_file/pprof/cpu-H.prof "$pprof_url/profile?seconds=180"
curl -o $pprof_file/pprof/heap-H.prof "$pprof_url/heap"
curl -o $pprof_file/pprof/goroutine-H.prof "$pprof_url/goroutine"

# 生成 火焰圖
go tool pprof -svg -output=$pprof_file/svg/cpu-H.svg $pprof_file/pprof/cpu-H.prof
go tool pprof -svg -output=$pprof_file/svg/heap-H.svg $pprof_file/pprof/heap-H.prof
go tool pprof -svg -output=$pprof_file/svg/goroutine-H.svg $pprof_file/pprof/goroutine-H.prof
