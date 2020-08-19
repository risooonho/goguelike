// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tilearea

import (
	"fmt"
	"image"
	"image/color"

	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/tileoptype"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/wrapper"
)

type TileArea [][]tile_flag.TileFlag

func (ta TileArea) String() string {
	return fmt.Sprintf("TileArea[%v %v]", len(ta), len(ta[0]))
}

func New(x, y int) TileArea {
	ta := make(TileArea, x)
	for i, _ := range ta {
		ta[i] = make([]tile_flag.TileFlag, y)
	}
	return ta
}

func (ta TileArea) Dup() TileArea {
	w, h := len(ta), len(ta[0])
	rtn := make(TileArea, w)
	for x, _ := range rtn {
		rtn[x] = make([]tile_flag.TileFlag, h)
		copy(rtn[x], ta[x])
	}
	return rtn
}

func (ta TileArea) DupWithFilter(fn func(x, y int) bool) TileArea {
	w, h := len(ta), len(ta[0])
	rtn := make(TileArea, w)
	for x, xv := range ta {
		rtn[x] = make([]tile_flag.TileFlag, h)
		for y, yv := range xv {
			if fn(x, y) {
				rtn[x][y] = yv
			}
		}
	}
	return rtn
}

// for draw2d
func (ta TileArea) GetXYLen() (int, int) {
	return len(ta), len(ta[0])
}

// for draw2d
func (ta TileArea) OpXY(x, y int, v tile_flag.TileTypeValue) {
	// rv := v.(tile_flag.TileTypeValue)
	ta[x][y].Op(v)
}

func (ta TileArea) TotalPos() int {
	return len(ta) * len(ta[0])
}

func (ta TileArea) GetByXY(x, y int) *tile_flag.TileFlag {
	return &ta[x][y]
}

func (ta TileArea) DrawRoomsToFloor(rs []*room.Room) {
	xWrap, yWrap := ta.GetXYWrapper()
	for _, r := range rs {
		roomRect := r.Area
		for x, xv := range r.Tiles {
			for y, yv := range xv {
				tax, tay := xWrap(roomRect.X+x), yWrap(roomRect.Y+y)
				ta[tax][tay].Op(tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: r.BgTile})
				ta[tax][tay].Op(tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: yv})
			}
		}
	}
}

func (ta TileArea) GetXYWrapper() (func(int) int, func(int) int) {
	return wrapper.New(len(ta)).Wrap, wrapper.New(len(ta[0])).Wrap
}

func (ta TileArea) ToImage(zoom int) *image.RGBA {
	Xlen, Ylen := ta.GetXYLen()
	img := image.NewRGBA(image.Rect(0, 0, Xlen*zoom, Ylen*zoom))

	for srcY := 0; srcY < Ylen; srcY++ {
		for srcX := 0; srcX < Xlen; srcX++ {
			r, g, b := ta[srcX][srcY].ToRGB()
			co := color.RGBA{r, g, b, 0xff}
			for y := srcY * zoom; y < srcY*zoom+zoom; y++ {
				for x := srcX * zoom; x < srcX*zoom+zoom; x++ {
					img.Set(x, y, co)
				}
			}
		}
	}
	return img
}

func (ta TileArea) GetTiles() TileArea {
	return ta
}
