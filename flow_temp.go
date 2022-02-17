package main

import (
	"fmt"
	"os"
)

//网络流测试

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

const INF = int(1e18)

type Flow struct {
	st, ed             int
	maxFlow            int //最大匹配数
	head, nt, to, w, d []int
	cnt                int
	matchL, matchR     []int
	q                  []int
	h, t               int
	idx                int
}

func (f *Flow) add(x, y, z int) {
	f.cnt++
	f.nt[f.cnt] = f.head[x]
	f.head[x] = f.cnt
	f.to[f.cnt] = y
	f.w[f.cnt] = z
}
func (f *Flow) init(n, m, idx int) {
	f.cnt = 1
	f.head = make([]int, n)
	f.matchL = make([]int, n)
	f.matchR = make([]int, n)
	f.d = make([]int, n)
	f.q = make([]int, n)

	f.nt = make([]int, m)
	f.to = make([]int, m)
	f.w = make([]int, m)
	f.idx = idx
	f.idx++
	f.st = idx
	f.idx++
	f.ed = idx
}
func (f *Flow) bfs() bool {
	f.h = 0
	f.t = 0
	f.q[f.t] = f.st
	f.t++
	for i := int(0); i <= f.idx; i++ {
		f.d[i] = 0
	}
	f.d[f.st] = 1
	for f.h < f.t {
		x := f.q[f.h]
		f.h++
		for i := f.head[x]; i > 0; i = f.nt[i] {
			v := f.to[i]
			if f.w[i] != 0 && f.d[v] == 0 {
				f.d[v] = f.d[x] + 1
				if v == f.ed {
					return true
				}
				f.q[f.t] = v
				f.t++
			}
		}
	}
	return f.d[f.ed] != 0
}
func (f *Flow) dfs(x, flow int) int {
	if x == f.ed {
		return flow
	}
	res := flow
	for i := f.head[x]; i > 0; i = f.nt[i] {
		v := f.to[i]
		if f.w[i] > 0 && f.d[v] == f.d[x]+1 {
			k := f.dfs(v, Min(res, f.w[i]))
			f.w[i] -= k
			f.w[i^1] += k
			res -= k
			if k == 0 {
				f.d[v] = -1
			}
			if res == 0 {
				break
			}
		}
	}
	return flow - res
}
func (f *Flow) cal() {
	/*
		遍历残量网络获得方案
		正边是偶数索引，回边是奇数
	*/
	cnt := f.idx - 2 //去掉源点和汇点
	for x := int(1); x <= cnt; x++ {
		for i := f.head[x]; i > 0; i = f.nt[i] {
			v := f.to[i]
			if i%2 == 1 && f.w[i] > 0 { //回边
				f.matchR[x] = v
				f.matchL[v] = x
			}
		}
	}
}
func (f *Flow) run() {
	f.maxFlow = 0
	for f.bfs() {
		//f.dfs(f.st, INF)
		f.maxFlow += f.dfs(f.st, INF)
	}
	f.cal()
}

func main2() {
	f := &Flow{}
	var n, m, s, t int
	fmt.Fscan(os.Stdin, &n, &m, &s, &t)
	f.init(n+5, m*2+5, n)
	f.st = s
	f.ed = t
	for i := int(0); i < m; i++ {
		var a, b, c int
		fmt.Fscan(os.Stdin, &a, &b, &c)
		f.add(a, b, c)
		f.add(b, a, 0)
	}
	f.run()
	fmt.Println(f.maxFlow)
}
