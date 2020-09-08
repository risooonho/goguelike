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

package noroomtower

import "github.com/kasworld/goguelike/tool/towermaker/towermake"

func New(name string) *towermake.Tower {
	tw := towermake.New(name)
	tw.Add("SightTest", 64, 64, 64, 0, 1.0).Appends(
		"ResourceFillRect resource=Soil amount=1  x=0 y=0  w=64 h=64",
	)

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}

	for _, fm := range tw.GetList() {
		fm.Appendf(
			"AddAreaAttack%v display=LightHouse acttype=LightHouse degree=0 perturn=10 count=%v message=LightHouse",
			"Rand", 1,
		)
		fm.Appendf(
			"AddAreaAttack%v display=GateKeeper acttype=GateKeeper degree=0 perturn=10 count=%v message=GateKeeper",
			"Rand", 1,
		)
		fm.Appendf(
			"AddMine%v display=None count=%v message=Mine",
			"Rand", 2,
		)
	}
	return tw
}
