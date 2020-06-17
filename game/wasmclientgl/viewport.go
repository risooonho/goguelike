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

	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/lib/imagecanvas"
	"github.com/kasworld/wrapper"
)

type Viewport struct {
	rnd *rand.Rand

	ViewWidth  int
	ViewHeight int

	// for tile draw
	TileImgCnvList          [tile.Tile_Count]*imagecanvas.ImageCanvas
	textureTileWrapInfoList [tile.Tile_Count]textureTileWrapInfo

	DarkerTileImgCnv *imagecanvas.ImageCanvas

	CanvasGL      js.Value
	threejs       js.Value
	camera        js.Value
	light         js.Value
	renderer      js.Value
	fontLoader    js.Value
	textureLoader js.Value

	scene       js.Value
	jsSceneObjs map[g2id.G2ID]js.Value

	// title
	font_helvetiker_regular js.Value
	jsoTitle                js.Value

	// terrain
	floorG2ID2ClientField map[g2id.G2ID]*ClientField

	// cache
	colorMaterialCache map[uint32]js.Value

	aoGeometryCache map[factiontype.FactionType]js.Value
	foGeometryCache map[fieldobjdisplaytype.FieldObjDisplayType]js.Value

	eqGeometryCache     map[equipslottype.EquipSlotType]js.Value
	potionGeometryCache map[potiontype.PotionType]js.Value
	scrollGeometryCache map[scrolltype.ScrollType]js.Value
	moneyGeometry       js.Value
}

func NewViewport() *Viewport {
	vp := &Viewport{
		rnd:                   rand.New(rand.NewSource(time.Now().UnixNano())),
		jsSceneObjs:           make(map[g2id.G2ID]js.Value),
		floorG2ID2ClientField: make(map[g2id.G2ID]*ClientField),
		aoGeometryCache:       make(map[factiontype.FactionType]js.Value),
		colorMaterialCache:    make(map[uint32]js.Value),
	}

	vp.threejs = js.Global().Get("THREE")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.CanvasGL = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvasglholder").Call("appendChild", vp.CanvasGL)
	vp.CanvasGL.Set("tabindex", "1")

	vp.scene = vp.ThreeJsNew("Scene")
	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 60, 1, 1, HelperSize*2)
	vp.textureLoader = vp.ThreeJsNew("TextureLoader")
	vp.fontLoader = vp.ThreeJsNew("FontLoader")

	// for tile draw
	for i, v := range tile.TileScrollAttrib {
		if v.Texture {
			idstr := fmt.Sprintf("%vPng", tile.Tile(i))
			vp.TileImgCnvList[i] = imagecanvas.NewByID(idstr)
			vp.textureTileWrapInfoList[i] = textureTileWrapInfo{
				Xcount: vp.TileImgCnvList[i].W / CellSize,
				Ycount: vp.TileImgCnvList[i].H / CellSize,
				WrapX:  wrapper.New(vp.TileImgCnvList[i].W - CellSize).WrapSafe,
				WrapY:  wrapper.New(vp.TileImgCnvList[i].H - CellSize).WrapSafe,
			}
		}
	}
	vp.DarkerTileImgCnv = imagecanvas.NewByID("DarkerPng")

	vp.initHelpers()
	vp.initTitle()
	return vp
}

func (vp *Viewport) Hide() {
	vp.CanvasGL.Get("style").Set("display", "none")
}
func (vp *Viewport) Show() {
	vp.CanvasGL.Get("style").Set("display", "initial")
}

func (vp *Viewport) ResizeCanvas(title bool) {
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Int()
	winH := win.Get("innerHeight").Int()
	if title {
		winH /= 3
	}
	vp.CanvasGL.Call("setAttribute", "width", winW)
	vp.CanvasGL.Call("setAttribute", "height", winH)
	vp.ViewWidth = winW
	vp.ViewHeight = winH

	vp.camera.Set("aspect", float64(winW)/float64(winH))
	vp.camera.Call("updateProjectionMatrix")

	vp.CanvasGL.Call("setAttribute", "width", winW)
	vp.CanvasGL.Call("setAttribute", "height", winH)
	vp.renderer.Call("setSize", winW, winH)
}

func (vp *Viewport) Focus() {
	vp.CanvasGL.Call("focus")
}

func (vp *Viewport) Zoom(state int) {
}

func (vp *Viewport) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.CanvasGL.Call("addEventListener", evt, js.FuncOf(fn))
}

func (vp *Viewport) Draw() {
	vp.renderer.Call("render", vp.scene, vp.camera)
}

func (vp *Viewport) ThreeJsNew(name string, args ...interface{}) js.Value {
	return vp.threejs.Get(name).New(args...)
}
