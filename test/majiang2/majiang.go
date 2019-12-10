package main

import (
	"flag"
	"fmt"
	"github.com/fwhappy/mahjong/card"
	"github.com/fwhappy/mahjong/suggest"
	"github.com/fwhappy/mahjong/wall"
	"github.com/fwhappy/mahjong/win"
	"sort"
)

var (
	tiles = []int{1,1,1,1,2,2,2,2,3,3,3,3,4,4,4,4,5,5,5,5,6,6,6,6,7,7,7,7,8,8,8,8,9,9,9,9,11,11,11,11,12,12,12,12,13,13,13,13,14,14,14,14,15,15,15,15,16,16,16,16,17,17,17,17,18,18,18,18,19,19,19,19,41,41,41,41,42,42,42,42,43,43,43,43}
	majongMap = map[int]string {
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
		41: "红中",
		42: "发财",
		43: "白板",
	}

	p int
	currentIndex = 0
	handSlice = make([][]int, 3)
	showSlice = make([][]int, 3)
	discardSlice = make([][]int, 3)
	suggestSlice = make([]*suggest.MSelector, 3)
)

func init()  {
	flag.IntVar(&p, "p", 3, "参与游戏人数")
}

func main()  {
	// 三个 AI
	suggestSlice[0] = suggest.NewMSelector()
	suggestSlice[1] = suggest.NewMSelector()
	suggestSlice[2] = suggest.NewMSelector()
	suggestSlice[0].SetAILevel(suggest.AI_KING)
	suggestSlice[1].SetAILevel(suggest.AI_BRASS)
	suggestSlice[2].SetAILevel(suggest.AI_PLATINUM)

	// 初始化牌墙
	w := wall.NewWall()
	w.SetTiles(tiles)
	// 洗牌
	w.Shuffle()
	// 初始化手牌
	initHand(w)
    // 打印手牌
	for k,v := range handSlice {
		fmt.Print(k+1, "号玩家", "\t")
		sort.Ints(v)
		outPutWall(v)
		fmt.Println("")
	}

	// 打牌阶段
	for w.IsAllDrawn() {
		// 是否蹦，杠
		card.GetRelationTiles()
		// 揭牌
		handSlice[currentIndex] = append(handSlice[currentIndex], w.ForwardDraw())
		// 判断是否胡牌
		if canWin := win.CanWin(handSlice[currentIndex], showSlice[currentIndex]); canWin {
			fmt.Println(currentIndex+1, "号玩家获胜")
		}
		// 设置ai
		suggestSlice[currentIndex].SetHandTilesSlice(handSlice[currentIndex])
		suggestSlice[currentIndex].SetShowTilesSlice(showSlice[currentIndex])
		suggestSlice[currentIndex].SetDiscardTilesSlice(discardSlice[currentIndex])
		tile := suggestSlice[currentIndex].GetSuggest()
		// 打牌
		delTile(currentIndex, tile)
		// 是否听牌
		// ting.CanTing(handSlice[currentIndex], showSlice[currentIndex])
	}
}

func delTile(currentIndex int, tile int) {
	for i:=0; i<=len(handSlice[currentIndex]); i++ {
		if handSlice[currentIndex][i] == tile {
			handSlice[currentIndex] = append(handSlice[currentIndex][:i], handSlice[currentIndex][i+1:]...)
			break
		}
	}
	discardSlice[currentIndex] = append(discardSlice[currentIndex], tile)
}

func outPutWall(tiles []int)  {
	for _,v :=range tiles {
		fmt.Print(majongMap[v], "\t")
	}
}

func initHand(w *wall.Wall)  {
	fmt.Println(p, "人游戏")
	// 第一次抓牌，4张
	handSlice[0] = w.ForwardDrawMulti(4)
	handSlice[1] = w.ForwardDrawMulti(4)
	handSlice[2] = w.ForwardDrawMulti(4)
	// 第二次抓牌，4张
	handSlice[0] = append(handSlice[0], w.ForwardDrawMulti(4)...)
	handSlice[1] = append(handSlice[1], w.ForwardDrawMulti(4)...)
	handSlice[2] = append(handSlice[2], w.ForwardDrawMulti(4)...)
	// 第三次抓牌，4张
	handSlice[0] = append(handSlice[0], w.ForwardDrawMulti(4)...)
	handSlice[1] = append(handSlice[1], w.ForwardDrawMulti(4)...)
	handSlice[2] = append(handSlice[2], w.ForwardDrawMulti(4)...)
	// 第四次抓牌
	handSlice[0] = append(handSlice[0], w.ForwardDraw())
	handSlice[1] = append(handSlice[1], w.ForwardDraw())
	handSlice[2] = append(handSlice[2], w.ForwardDraw())
}
