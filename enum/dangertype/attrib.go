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

package dangertype

import (
	"github.com/kasworld/htmlcolors"
)

func (dt DangerType) Turn2Live() int {
	return attrib[dt].turn2Live
}
func (dt DangerType) Color24() htmlcolors.Color24 {
	return attrib[dt].color
}

func (dt DangerType) Scale4UI() float64 {
	return attrib[dt].scale4ui
}

var attrib = [DangerType_Count]struct {
	turn2Live int
	scale4ui  float64
	color     htmlcolors.Color24
}{
	None:             {0, 1.0, htmlcolors.Black},
	BasicAttack:      {1, 1.0, htmlcolors.Red},
	WideAttack:       {1, 1.0, htmlcolors.Crimson},
	LongAttack:       {1, 1.0, htmlcolors.FireBrick},
	RotateLineAttack: {1, 2.0, htmlcolors.DeepPink},
	MineExplode:      {1, 2.0, htmlcolors.Orange},
}
