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

package wasmclientgl

import (
	"fmt"
	"syscall/js"

	"github.com/kasworld/goguelike/game/clientinitdata"
	"github.com/kasworld/goguelike/game/soundmap"
	"github.com/kasworld/goguelike/lib/htmlbutton"
	"github.com/kasworld/gowasmlib/jslog"
)

var gameOptions *htmlbutton.HTMLButtonGroup

// prevent compiler initialize loop error
var _gameopt = htmlbutton.NewButtonGroup("Options",
	[]*htmlbutton.HTMLButton{
		htmlbutton.New("q", "LeftInfo", []string{"LeftInfoOff", "LeftInfoOn"},
			"show/hide left info", cmdLeftInfo, 1),
		htmlbutton.New("w", "CenterInfo", []string{"HelpOff", "Highscore", "ClientInfo", "Help", "FactionInfo",
			"CarryObjectInfo", "PotionInfo", "ScrollInfo", "MoneyColor",
			"TileInfo", "ConditionInfo", "FieldObjInfo"},
			"rotate help info", cmdRotateCenterInfo, 0),
		htmlbutton.New("e", "RightInfo", []string{
			"RightInfoOff", "Message", "DebugInfo", "InvenList", "FieldObjList", "FloorList"},
			"Rotate right info", cmdRotateRightInfo, 1),
		htmlbutton.New("r", "ViewMode", []string{"PlayVP", "FloorVP"},
			"play view / floor view", cmdToggleVPFloorPlay, 0),
		htmlbutton.New("t", "Zoom", []string{"Zoom0", "Zoom1", "Zoom2"},
			"Zoom viewport", cmdToggleZoom, 0),
		htmlbutton.New("y", "Angle", []string{"Angle0", "Angle1", "Angle2"},
			"Angle viewport", cmdToggleAngle, 0),
		htmlbutton.New("u", "Sound", []string{"SoundOn", "SoundOff"},
			"Sound on/off", cmdToggleSound, 1),
	})

func cmdLeftInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}

func (app *WasmClient) updateLeftInfo() {
	v := gameOptions.GetByIDBase("LeftInfo")
	infoobj := GetElementById("leftinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // leftinfo
		infoobj.Set("innerHTML",
			app.makeBaseInfoHTML()+app.makeBuffListHTML())
	}
}

func cmdRotateCenterInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.updateCenterInfo()
	v.Blur()
}

func (app *WasmClient) updateCenterInfo() {
	v := gameOptions.GetByIDBase("CenterInfo")
	infoobj := GetElementById("centerinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", app.MakeNotiMessage())
	case 1: // highscore
		go func() {
			infoobj.Set("innerHTML", clientinitdata.LoadHighScoreHTML())
		}()
	case 2: // clientinfo
		infoobj.Set("innerHTML", clientinitdata.MakeClientInfoHTML())
	case 3: // helpinfo
		infoobj.Set("innerHTML", MakeHelpInfoHTML())
	case 4: // faction
		infoobj.Set("innerHTML", clientinitdata.MakeHelpFactionHTML())
	case 5: // carryobj
		infoobj.Set("innerHTML", clientinitdata.MakeHelpCarryObjectHTML())
	case 6: // potion
		infoobj.Set("innerHTML", clientinitdata.MakeHelpPotionHTML())
	case 7: // scroll
		infoobj.Set("innerHTML", clientinitdata.MakeHelpScrollHTML())
	case 8: // Money color
		infoobj.Set("innerHTML", clientinitdata.MakeHelpMoneyColorHTML())
	case 9: // tile
		infoobj.Set("innerHTML", clientinitdata.MakeHelpTileHTML())
	case 10: // condition
		infoobj.Set("innerHTML", clientinitdata.MakeHelpConditionHTML())
	case 11: // fieldobj
		infoobj.Set("innerHTML", clientinitdata.MakeHelpFieldObjHTML())
	}
	winW := js.Global().Get("window").Get("innerWidth").Float()
	infoobjW := infoobj.Get("offsetWidth").Float()
	GetElementById("centerinfo").Get("style").Set("left", fmt.Sprintf("%vpx", winW/2-infoobjW/2))

}

func cmdRotateRightInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.updateRightInfo()
	v.Blur()
}

func (app *WasmClient) updateRightInfo() {
	v := gameOptions.GetByIDBase("RightInfo")
	infoobj := GetElementById("rightinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // Message
		app.systemMessage = app.systemMessage.GetLastN(DisplayLineLimit)
		infoobj.Set("innerHTML", app.systemMessage.ToHtmlStringRev())
	case 2: // DebugInfo
		infoobj.Set("innerHTML", app.makeDebugInfoHTML())
	case 3: // InvenList
		infoobj.Set("innerHTML", app.makeInvenInfoHTML())
	case 4: // FieldObjList
		infoobj.Set("innerHTML", app.makeFieldObjListHTML())
	case 5: // FloorList
		infoobj.Set("innerHTML", app.makeFloorListHTML())
	}
}

func cmdToggleVPFloorPlay(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	switch v.State {
	case 0: // play viewpot mode
		// need rebuild tile, fieldobj
		cf := app.currentFloor()
		vpx := app.taNotiData.VPX
		vpy := app.taNotiData.VPY
		app.vp.makeClientTile4FloorView(cf, vpx, vpy)
		app.vp.updateFieldObjInView(cf, vpx, vpy)

	case 1: // floor viewport mode
		app.floorVPPosX, app.floorVPPosY = app.GetPlayerXY()
	}
	v.Blur()
}

func cmdToggleZoom(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}

	app.vp.Zoom(v.State)
	app.NotiMessage.AppendTf(tcsInfo, "Zoom%v", v.State)
	v.Blur()
}

func cmdToggleAngle(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}

	// if cf := app.currentFloor(); cf != nil {
	// 	cf.Angle(v.State)
	// }

	app.systemMessage.Appendf("Angle%v", v.State)
	v.Blur()
}

func cmdToggleSound(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	if v.State == 0 {
		soundmap.SoundOn = true
		app.systemMessage.Append("SoundOn")
		app.NotiMessage.AppendTf(tcsInfo, "SoundOn")
	} else {
		soundmap.SoundOn = false
		app.systemMessage.Append("SoundOff")
		app.NotiMessage.AppendTf(tcsInfo, "SoundOff")
	}
	v.Blur()
}

func (app *WasmClient) reset2Default() {
	btn := gameOptions.GetByIDBase("LeftInfo")
	btn.State = 1
	btn.UpdateButtonText()

	btn = gameOptions.GetByIDBase("CenterInfo")
	btn.State = 0
	btn.UpdateButtonText()

	btn = gameOptions.GetByIDBase("RightInfo")
	btn.State = 1
	btn.UpdateButtonText()

	btn = gameOptions.GetByIDBase("ViewMode")
	btn.State = 0
	btn.UpdateButtonText()

	btn = gameOptions.GetByIDBase("Zoom")
	btn.State = 1
	btn.UpdateButtonText()
	cmdToggleZoom(app, btn)

	btn = gameOptions.GetByIDBase("Angle")
	btn.State = 1
	btn.UpdateButtonText()
	cmdToggleAngle(app, btn)

	btn = gameOptions.GetByIDBase("Sound")
	btn.State = 1
	btn.UpdateButtonText()
	cmdToggleSound(app, btn)
}
