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
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/lib/htmlbutton"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wrapspan"
)

var commandButtons = htmlbutton.NewButtonGroup("Commands",
	[]*htmlbutton.HTMLButton{
		htmlbutton.New("a", "KillSelf", []string{"KillSelf"}, "Kill self", cmdKillSelf, 0),
		htmlbutton.New("s", "EnterPortal", []string{"EnterPortal"}, "Enter portal", cmdEnterPortal, 0),
		htmlbutton.New("d", "Teleport", []string{"Teleport"}, "Teleport random in floor", cmdTeleport, 0),
		htmlbutton.New("f", "Rebirth", []string{"Rebirth"}, "Rebirth", cmdRebirth, 0),
		htmlbutton.New("g", "ShowAchieve", []string{"ShowAchieve"}, "Show Achievement", cmdShowAchieve, 0),
		htmlbutton.New("h", "ShowPotionStat", []string{"ShowPotionStat"}, "Show PotionStat", cmdShowPotionStat, 0),
		htmlbutton.New("j", "ShowScrollStat", []string{"ShowScrollStat"}, "Show ScrollStat", cmdShowScrollStat, 0),
		htmlbutton.New("k", "ShowFOActStat", []string{"ShowFOActStat"}, "Show FOActStat", cmdShowFOActStat, 0),
		htmlbutton.New("l", "ShowActionStat", []string{"ShowActionStat"}, "Show ActionStat", cmdShowActionStat, 0),
		htmlbutton.New(";", "ShowConditionStat", []string{"ShowConditionStat"}, "Show ConditionStat", cmdShowConditionStat, 0),
	})

func cmdKillSelf(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.KillSelf,
		&c2t_obj.ReqKillSelf_data{},
	)
	v.Blur()
}

func cmdShowAchieve(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== Achievement == "))
			for i, v := range rpk.AchieveStat {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", achievetype.AchieveType(i).String(), v))
			}
			return nil
		},
	)
	v.Blur()
}

func cmdShowPotionStat(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== Portion stat == "))
			for i, v := range rpk.PotionStat {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", potiontype.PotionType(i).String(), v))
			}
			return nil
		},
	)
	v.Blur()
}
func cmdShowScrollStat(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== Scroll stat == "))
			for i, v := range rpk.ScrollStat {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", scrolltype.ScrollType(i).String(), v))
			}
			return nil
		},
	)
	v.Blur()
}
func cmdShowFOActStat(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== FieldObj act stat == "))
			for i, v := range rpk.FOActStat {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", fieldobjacttype.FieldObjActType(i).String(), v))
			}
			return nil
		},
	)
	v.Blur()
}
func cmdShowActionStat(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== Action stat == "))
			for i, v := range rpk.AOActionStat {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", c2t_idcmd.CommandID(i).String(), v))
			}
			return nil
		},
	)
	v.Blur()
}
func cmdShowConditionStat(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== Condition stat == "))
			for i, v := range rpk.ConditionStat {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", condition.Condition(i).String(), v))
			}
			return nil
		},
	)
	v.Blur()
}

func cmdEnterPortal(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.EnterPortal,
		&c2t_obj.ReqEnterPortal_data{},
	)
	v.Blur()
}

func cmdTeleport(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.ActTeleport,
		&c2t_obj.ReqActTeleport_data{},
	)
	v.Blur()
}

func cmdRebirth(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	if app.olNotiData.ActiveObj.HP <= 0 {
		go app.sendPacket(c2t_idcmd.Rebirth,
			&c2t_obj.ReqRebirth_data{},
		)
	} else {
		app.systemMessage.Append(
			wrapspan.ColorText("OrangeRed",
				"no need to rebirth"))
	}
	v.Blur()
}
