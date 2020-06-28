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
	"syscall/js"
)

type Viewport struct {
	ViewWidth  int
	ViewHeight int

	CanvasGL js.Value
	renderer js.Value
}

func NewViewport() *Viewport {
	vp := &Viewport{}

	vp.renderer = ThreeJsNew("WebGLRenderer")
	vp.CanvasGL = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvasglholder").Call("appendChild", vp.CanvasGL)
	vp.CanvasGL.Set("tabindex", "1")

	return vp
}

func (vp *Viewport) Hide() {
	vp.CanvasGL.Get("style").Set("display", "none")
}
func (vp *Viewport) Show() {
	vp.CanvasGL.Get("style").Set("display", "initial")
}

func (vp *Viewport) ResizeCanvas(title bool, cf *ClientFloorGL) {
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Int()
	winH := win.Get("innerHeight").Int()
	if title {
		winH /= 3
	}
	vp.ViewWidth = winW
	vp.ViewHeight = winH

	vp.CanvasGL.Call("setAttribute", "width", winW)
	vp.CanvasGL.Call("setAttribute", "height", winH)
	vp.renderer.Call("setSize", winW, winH)

	if cf != nil {
		cf.camera.Set("aspect", float64(winW)/float64(winH))
		cf.camera.Call("updateProjectionMatrix")
	}
}

func (vp *Viewport) Focus() {
	vp.CanvasGL.Call("focus")
}

func (vp *Viewport) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.CanvasGL.Call("addEventListener", evt, js.FuncOf(fn))
}
