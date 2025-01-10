#!/bin/bash

# 設定 pprof 伺服器的地址
pprof_url="http://localhost:6060/debug/pprof"

# 生成各種性能報告
curl -o /Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/pprof-Output/pprof/cpu.prof "$pprof_url/profile?seconds=300"
curl -o /Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/pprof-Output/pprof/heap.prof "$pprof_url/heap"
curl -o /Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/pprof-Output/pprof/goroutine.prof "$pprof_url/goroutine"

echo "All pprof data has been saved."

# 生成火焰圖（SVG）
echo "Generating flamegraphs (SVG)..."

# 生成 CPU 火焰圖
go tool pprof -svg -output=/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/pprof-Output/svg/cpu.svg cpu.prof

# 生成 Heap 火焰圖
go tool pprof -svg -output=/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/pprof-Output/svg/heap.svg heap.prof

# 生成 Goroutine 火焰圖
go tool pprof -svg -output=/Users/cyrusman/Desktop/ProgrammingLearning/Golang-Blockchain-2024/blockchain-back/pprof-Output/svg/goroutine.svg goroutine.prof

echo "All flamegraphs have been saved as SVG."