#!/bin/bash

# 設定 pprof 伺服器的地址
pprof_url="http://localhost:6060/debug/pprof"
pprof_file="/Users/cyrusman/Desktop/ProgrammingLearning/Project/Golang-Blockchain-2024/blockchain-back/pprof-Output"

# 生成各種性能報告
curl -o $pprof_file/pprof/cpu-L.prof "$pprof_url/profile?seconds=180"
curl -o $pprof_file/pprof/heap-L.prof "$pprof_url/heap"
curl -o $pprof_file/pprof/goroutine-L.prof "$pprof_url/goroutine"

# 生成 火焰圖
go tool pprof -svg -output=$pprof_file/svg/cpu-L.svg $pprof_file/pprof/cpu-L.prof
go tool pprof -svg -output=$pprof_file/svg/heap-L.svg $pprof_file/pprof/heap-L.prof
go tool pprof -svg -output=$pprof_file/svg/goroutine-L.svg $pprof_file/pprof/goroutine-L.prof
