package knife

import (
	"fmt"
	"image-server/pkg/objects"
	"testing"
)

func TestComposite(t *testing.T) {
	_, err := Composite([]CompositeParam{
		{
			Text: "笔船一日游",
			X:    370,
			Y:    900,
			Font: &objects.Font{
				Size:  30,
				Color: "#333333",
				Align: "center",
			},
		},
		{
			Url:    "https://qn-album.meetb.cn/506b4202008181936506586.jpg",
			X:      309,
			Y:      669,
			Width:  132,
			Height: 132,
		},
		{
			Url:    "https://node.meetb.cn/img/posters/project_detail.png",
			X:      0,
			Y:      0,
			Width:  750,
			Height: 1334,
		},
		{
			Text: "姚木木",
			X:    335,
			Y:    840,
			Font: &objects.Font{
				Size:  30,
				Color: "#333333",
				Align: "center",
			},
		},
		{
			Text: "￥999",
			X:    25,
			Y:    615,
			Font: &objects.Font{
				Size:  42,
				Color: "#ffffff",
				Align: "left",
			},
		},
		{
			Text: "2020/09/08",
			X:    30,
			Y:    700,
			Font: &objects.Font{
				Size:  26,
				Color: "#ffffff",
				Align: "left",
			},
		},
		{
			Text: "绿地中心",
			X:    610,
			Y:    700,
			Font: &objects.Font{
				Size:  26,
				Color: "#ffffff",
				Align: "right",
			},
		},
		{
			Url:    "https://source.meetb.cn/inst_travel/posters/share_program/64/4856a0b28c9ee4fbbb3d595059576846(140x140)_qr.png",
			X:      258,
			Y:      960,
			Width:  234,
			Height: 234,
		},
		{
			Url:    "https://qn-album.meetb.cn/506b4202008181936506586.jpg",
			X:      600,
			Y:      300,
			Width:  200,
			Height: 200,
			Circle: true,
		},
	})
	if err != nil {
		panic(fmt.Errorf("生成失败 %s", err.Error()))
	}
}
