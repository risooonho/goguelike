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
	"sync"
	"syscall/js"

	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func NewCarryObj3DGeo(str string) js.Value {
	refSize := float64(DstCellSize) / 2
	geo := ThreeJsNew("TextGeometry", str,
		map[string]interface{}{
			"font":           gFont_helvetiker_regular,
			"size":           refSize * 0.7,
			"height":         refSize * 0.3,
			"curveSegments":  refSize / 3,
			"bevelEnabled":   true,
			"bevelThickness": refSize / 8,
			"bevelSize":      refSize / 16,
			"bevelOffset":    0,
			"bevelSegments":  refSize / 8,
		})
	geo.Call("center")
	return geo
}

func Equiped2StrColor(o *c2t_obj.EquipClient) (string, string) {
	return o.EquipType.Rune(), o.Faction.Color24().ToHTMLColorString()
}

func CarryObj2StrColor(o *c2t_obj.CarryObjClientOnFloor) (string, string) {
	switch o.CarryingObjectType {
	case carryingobjecttype.Equip:
		return o.EquipType.Rune(), o.Faction.Color24().ToHTMLColorString()
	case carryingobjecttype.Money:
		for _, v := range moneycolor.Attrib {
			if o.Value < v.UpLimit {
				return "$", v.Color.ToHTMLColorString()
			}
		}
		v := moneycolor.Attrib[len(moneycolor.Attrib)-1]
		return "$", v.Color.ToHTMLColorString()
	case carryingobjecttype.Potion:
		return o.PotionType.Rune(), o.PotionType.Color24().ToHTMLColorString()
	case carryingobjecttype.Scroll:
		return o.ScrollType.Rune(), o.ScrollType.Color24().ToHTMLColorString()
	}
	return "", ""
}

var gPoolCarryObj3DGeo = NewPoolCarryObj3DGeo()

type PoolCarryObj3DGeo struct {
	mutex    sync.Mutex
	poolData map[string][]js.Value
	newCount int
	getCount int
	putCount int
}

func NewPoolCarryObj3DGeo() *PoolCarryObj3DGeo {
	return &PoolCarryObj3DGeo{
		poolData: make(map[string][]js.Value),
	}
}

func (p *PoolCarryObj3DGeo) String() string {
	return fmt.Sprintf("PoolCarryObj3DGeo[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolCarryObj3DGeo) Get(str string) js.Value {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn js.Value
	if l := len(p.poolData[str]); l > 0 {
		rtn = p.poolData[str][l-1]
		p.poolData[str] = p.poolData[str][:l-1]
		p.getCount++
	} else {
		rtn = NewCarryObj3DGeo(str)
		p.newCount++
	}
	return rtn
}

func (p *PoolCarryObj3DGeo) Put(geo js.Value, str string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[str] = append(p.poolData[str], geo)
	p.putCount++
}

type CarryObj3D struct {
	Str     string
	Color   string
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewCarryObj3D(str, color string) *CarryObj3D {
	mat := gPoolColorMaterial.Get((color))
	mat.Set("opacity", 1)
	geo := gPoolCarryObj3DGeo.Get(str)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &CarryObj3D{
		Str:     str,
		Color:   color,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func (aog *CarryObj3D) SetFieldPosition(fx, fy int, shInfo ShiftInfo) {
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+shInfo.X,
		-float64(fy)*DstCellSize-shInfo.Y,
		aog.GeoInfo.Len[2]/2+shInfo.Z,
	)
}

func (aog *CarryObj3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *CarryObj3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *CarryObj3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *CarryObj3D) Dispose() {
	// mesh do not need dispose
	// aog.Mesh.Get("geometry").Call("dispose")
	gPoolCarryObj3DGeo.Put(aog.Mesh.Get("geometry"), aog.Str)
	gPoolColorMaterial.Put(aog.Mesh.Get("material"))
	// aog.Mesh.Get("material").Call("dispose")
	aog.Mesh = js.Undefined()
	// no need createElement canvas dom obj
}

type ShiftInfo struct {
	X float64
	Y float64
	Z float64
}

// equipped shift, around ao
var aoEqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	// center
	equipslottype.Helmet:   {DstCellSize * 0.5, DstCellSize * 0.00, DstCellSize * 1.00},
	equipslottype.Amulet:   {DstCellSize * 0.5, DstCellSize * 0.33, DstCellSize * 1.00},
	equipslottype.Armor:    {DstCellSize * 0.5, DstCellSize * 0.66, DstCellSize * 1.00},
	equipslottype.Footwear: {DstCellSize * 0.5, DstCellSize * 1.00, DstCellSize * 1.00},

	// right
	equipslottype.Weapon:   {DstCellSize * 1.00, DstCellSize * 0.33, DstCellSize * 1.00},
	equipslottype.Gauntlet: {DstCellSize * 1.00, DstCellSize * 0.66, DstCellSize * 1.00},

	// left
	equipslottype.Shield: {DstCellSize * 0.00, DstCellSize * 0.33, DstCellSize * 1.00},
	equipslottype.Ring:   {DstCellSize * 0.00, DstCellSize * 0.66, DstCellSize * 1.00},
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

// on floor in tile
var eqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {DstCellSize * 0.0, DstCellSize * 0.0, DstCellSize * 0.33},
	equipslottype.Amulet: {DstCellSize * 0.75, DstCellSize * 0.0, DstCellSize * 0.33},

	equipslottype.Weapon: {DstCellSize * 0.0, DstCellSize * 0.25, DstCellSize * 0.33},
	equipslottype.Shield: {DstCellSize * 0.75, DstCellSize * 0.25, DstCellSize * 0.33},

	equipslottype.Ring:     {DstCellSize * 0.0, DstCellSize * 0.50, DstCellSize * 0.33},
	equipslottype.Gauntlet: {DstCellSize * 0.75, DstCellSize * 0.50, DstCellSize * 0.33},

	equipslottype.Armor:    {DstCellSize * 0.0, DstCellSize * 0.75, DstCellSize * 0.33},
	equipslottype.Footwear: {DstCellSize * 0.75, DstCellSize * 0.75, DstCellSize * 0.33},
}

var otherCarryObjShift = [carryingobjecttype.CarryingObjectType_Count]ShiftInfo{
	carryingobjecttype.Money:  {DstCellSize * 0.33, DstCellSize * 0.0, DstCellSize * 0.33},
	carryingobjecttype.Potion: {DstCellSize * 0.33, DstCellSize * 0.33, DstCellSize * 0.33},
	carryingobjecttype.Scroll: {DstCellSize * 0.33, DstCellSize * 0.66, DstCellSize * 0.33},
}
