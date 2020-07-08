// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wasmclientgl

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/kasworld/findnear"
	"github.com/kasworld/go-abs"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/clientinitdata"
	"github.com/kasworld/goguelike/lib/clienttile"
	"github.com/kasworld/goguelike/lib/webtilegroup"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

const (
	DisplayLineLimit = 3*gameconst.ViewPortH - gameconst.ViewPortH/2
	DstCellSize      = 32
	HelperSize       = DstCellSize * 32
	ClientViewLen    = 40
)

var gRnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

var gInitData *clientinitdata.InitData = clientinitdata.New()
var gClientTile *clienttile.ClientTile = clienttile.New()

var gXYLenListView = findnear.NewXYLenList(ClientViewLen, ClientViewLen)

var gTextureLoader js.Value = ThreeJsNew("TextureLoader")

var gColorMaterialCache map[string]js.Value = make(map[string]js.Value)

func GetColorMaterialByCache(co string) js.Value {
	mat, exist := gColorMaterialCache[co]
	if !exist {
		mat = ThreeJsNew("MeshPhongMaterial",
			map[string]interface{}{
				"color": co,
			},
		)
		mat.Set("transparent", true)
		gColorMaterialCache[co] = mat
	}
	return mat
}

func NewTileMaterial(ti webtilegroup.TileInfo) js.Value {
	Cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", DstCellSize)
	Cnv.Set("height", DstCellSize)
	Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	mat := ThreeJsNew("MeshPhongMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	mat.Set("transparent", true)
	// mat.Set("side", ThreeJs().Get("DoubleSide"))
	return mat
}

var gTileMaterialCache map[webtilegroup.TileInfo]js.Value = make(map[webtilegroup.TileInfo]js.Value)

func GetTileMaterialByCache(ti webtilegroup.TileInfo) js.Value {
	mat, exist := gTileMaterialCache[ti]
	if !exist {
		mat = NewTileMaterial(ti)
		gTileMaterialCache[ti] = mat
	}
	return mat
}

func NewTextureTileMaterial(ti tile.Tile) js.Value {
	Cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", DstCellSize)
	Cnv.Set("height", DstCellSize)

	img := GetElementById(fmt.Sprintf("%vPng", ti))
	Ctx.Call("drawImage", img,
		0, 0, DstCellSize, DstCellSize,
		0, 0, DstCellSize, DstCellSize)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	mat := ThreeJsNew("MeshPhongMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	mat.Set("transparent", true)
	// mat.Set("side", ThreeJs().Get("DoubleSide"))
	return mat
}

type textGeoKey struct {
	Str  string
	Size float64
}

var gFontLoader js.Value = ThreeJsNew("FontLoader")
var gFont_helvetiker_regular js.Value

var gTextGeometryCache map[textGeoKey]js.Value = make(map[textGeoKey]js.Value)

func GetTextGeometryByCache(str string, size float64) js.Value {
	geo, exist := gTextGeometryCache[textGeoKey{str, size}]
	curveSegments := size / 3
	if curveSegments < 1 {
		curveSegments = 1
	}
	bevelEnabled := true
	if size < 16 {
		bevelEnabled = false
	}
	bevelThickness := size / 8
	if bevelThickness < 1 {
		bevelThickness = 1
	}
	bevelSize := size / 16
	if bevelSize < 1 {
		bevelSize = 1
	}
	bevelSegments := size / 8
	if bevelSegments < 1 {
		bevelSegments = 1
	}
	if !exist {
		geo = ThreeJsNew("TextGeometry", str,
			map[string]interface{}{
				"font":           gFont_helvetiker_regular,
				"size":           size,
				"height":         5,
				"curveSegments":  curveSegments,
				"bevelEnabled":   bevelEnabled,
				"bevelThickness": bevelThickness,
				"bevelSize":      bevelSize,
				"bevelOffset":    0,
				"bevelSegments":  bevelSegments,
			})
		gTextGeometryCache[textGeoKey{str, size}] = geo
	}
	return geo
}

var gBoxGeometryCache map[[3]int]js.Value = make(map[[3]int]js.Value)

func GetBoxGeometryByCache(x, y, z int) js.Value {
	geo, exist := gBoxGeometryCache[[3]int{x, y, z}]
	if !exist {
		geo = ThreeJsNew("BoxGeometry", x, y, z)
		gBoxGeometryCache[[3]int{x, y, z}] = geo
	}
	return geo
}

func CalcCurrentFrame(difftick int64, fps float64) int {
	diffsec := float64(difftick) / float64(time.Second)
	frame := fps * diffsec
	return int(frame)
}

func CalcShiftDxDy(frameProgress float64) (int, int) {
	rate := 1 - frameProgress
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	dx := int(float64(DstCellSize) * rate)
	dy := int(float64(DstCellSize) * rate)
	return dx, dy
}

func SetPosition(jso js.Value, pos ...interface{}) {
	po := jso.Get("position")
	if len(pos) >= 1 {
		po.Set("x", pos[0])
	}
	if len(pos) >= 2 {
		po.Set("y", pos[1])
	}
	if len(pos) >= 3 {
		po.Set("z", pos[2])
	}
}

func ThreeJsNew(name string, args ...interface{}) js.Value {
	return js.Global().Get("THREE").Get(name).New(args...)
}

func ThreeJs() js.Value {
	return js.Global().Get("THREE")
}

func GetElementById(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func CarryObjClientOnFloor2DrawInfo(
	co *c2t_obj.CarryObjClientOnFloor) ShiftInfo {
	switch co.CarryingObjectType {
	default:
		return otherCarryObjShift[co.CarryingObjectType]
	case carryingobjecttype.Equip:
		return eqPosShift[co.EquipType]
	}
}

type ShiftInfo struct {
	X float64
	Y float64
	Z float64
}

// equipped shift, around ao
var aoEqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {-0.33, 0.0, 0.66},
	equipslottype.Amulet: {1.00, 0.0, 0.66},

	equipslottype.Weapon: {-0.33, 0.25, 0.66},
	equipslottype.Shield: {1.00, 0.25, 0.66},

	equipslottype.Ring:     {-0.33, 0.50, 0.66},
	equipslottype.Gauntlet: {1.00, 0.50, 0.66},

	equipslottype.Armor:    {-0.33, 0.75, 0.66},
	equipslottype.Footwear: {1.00, 0.75, 0.66},
}

// on floor in tile
var eqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {0.0, 0.0, 0.33},
	equipslottype.Amulet: {0.75, 0.0, 0.33},

	equipslottype.Weapon: {0.0, 0.25, 0.33},
	equipslottype.Shield: {0.75, 0.25, 0.33},

	equipslottype.Ring:     {0.0, 0.50, 0.33},
	equipslottype.Gauntlet: {0.75, 0.50, 0.33},

	equipslottype.Armor:    {0.0, 0.75, 0.33},
	equipslottype.Footwear: {0.75, 0.75, 0.33},
}

var otherCarryObjShift = [carryingobjecttype.CarryingObjectType_Count]ShiftInfo{
	carryingobjecttype.Money:  {0.33, 0.0, 0.33},
	carryingobjecttype.Potion: {0.33, 0.33, 0.33},
	carryingobjecttype.Scroll: {0.33, 0.66, 0.33},
}

// make fx,fy around vpx, vpy
func calcAroundPos(w, h, vpx, vpy, fx, fy int) (int, int) {
	if abs.Absi(fx-vpx) > w/2 {
		if fx > vpx {
			fx -= w
		} else {
			fx += w
		}
	}
	if abs.Absi(fy-vpy) > h/2 {
		if fy > vpy {
			fy -= h
		} else {
			fy += h
		}
	}
	return fx, fy
}

func makeEquipedMesh(o *c2t_obj.EquipClient) js.Value {
	ti := gClientTile.EquipTiles[o.EquipType][o.Faction]
	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(
		DstCellSize/3, DstCellSize/3, DstCellSize/3,
	)
	return ThreeJsNew("Mesh", geo, mat)
}

func makeCarryObjMesh(o *c2t_obj.CarryObjClientOnFloor) js.Value {
	var ti webtilegroup.TileInfo
	switch o.CarryingObjectType {
	case carryingobjecttype.Equip:
		ti = gClientTile.EquipTiles[o.EquipType][o.Faction]
	case carryingobjecttype.Money:
		var find bool
		for i, v := range moneycolor.Attrib {
			if o.Value < v.UpLimit {
				ti = gClientTile.GoldTiles[i]
				find = true
				break
			}
		}
		if !find {
			ti = gClientTile.GoldTiles[len(gClientTile.GoldTiles)-1]
		}
	case carryingobjecttype.Potion:
		ti = gClientTile.PotionTiles[o.PotionType]
	case carryingobjecttype.Scroll:
		ti = gClientTile.ScrollTiles[o.ScrollType]
	}
	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(
		DstCellSize/3, DstCellSize/3, DstCellSize/3,
	)
	return ThreeJsNew("Mesh", geo, mat)
}

func newFieldObjAt(o *c2t_obj.FieldObjClient, fx, fy int) js.Value {
	tlList := gClientTile.FieldObjTiles[o.DisplayType]
	tilediff := fx*5 + fy*3
	if tilediff < 0 {
		tilediff = -tilediff
	}
	ti := tlList[tilediff%len(tlList)]

	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(DstCellSize-1, DstCellSize-1, DstCellSize-1)
	geoInfo := GetGeoInfo(geo)
	mesh := ThreeJsNew("Mesh", geo, mat)
	SetPosition(
		mesh,
		float64(fx)*DstCellSize+geoInfo.Len[0]/2,
		-float64(fy)*DstCellSize-geoInfo.Len[1]/2,
		geoInfo.Len[2]/2)
	return mesh
}

func MakeFieldObjMatGeo(
	o *c2t_obj.FieldObjClient, fx, fy int) (js.Value, js.Value) {

	tlList := gClientTile.FieldObjTiles[o.DisplayType]
	tilediff := fx*5 + fy*3
	ti := tlList[tilediff%len(tlList)]
	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(
		DstCellSize-1, DstCellSize-1, DstCellSize-1)
	return mat, geo
}
