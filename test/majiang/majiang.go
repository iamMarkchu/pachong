package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
)

var (
	list []int
	majiangMap = map[int]string {
		1: "一筒",
		2: "二筒",
		3: "三筒",
		4: "四筒",
		5: "五筒",
		6: "六筒",
		7: "七筒",
		8: "八筒",
		9: "九筒",
		11: "一条",
		12: "二条",
		13: "三条",
		14: "四条",
		15: "五条",
		16: "六条",
		17: "七条",
		18: "八条",
		19: "九条",
		21: "红中",
		22: "发财",
		23: "白板",
	}
)

func init()  {
	for i :=1; i<=23; i++ {
		if i == 10 || i == 20 {
			continue
		}
		for j :=0; j<=3; j++ {
			list = append(list, i)
		}
	}
}
// 发13张牌，然后穷举， 13张牌之外的所有牌，是否满足规则 3n+2
func main()  {
	var currentList []int
	for i:=0; i<=13;i++ {
		if n, err := rand.Int(rand.Reader, big.NewInt(int64(len(list)))); err == nil {
			index := int(n.Int64())
			currentList = append(currentList, list[index])
			list = append(list[:index], list[index+1:]...)
		}
	}
	sort.Ints(currentList)
	for _,v := range currentList {
		fmt.Print(majiangMap[v])
	}
}
