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

package starttower

import (
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/tool/towermaker/floortemplate"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func New(name string) *towermake.Tower {
	var rnd = g2rand.New()

	tw := towermake.New(name)
	tw.Add("Practice", 64, 32, 0, 0, 0.7).Appends(
		floortemplate.Practice64x32()...,
	)
	tw.Add("SoilPlant", 64, 64, 16, 0, 1.0).Appends(
		floortemplate.SoilPlant64x64()...,
	)
	tw.Add("ManyPortals", 128, 128, 64, 0, 1.0).Appends(
		floortemplate.BGTile9Rooms128x128()...,
	)
	tw.Add("SoilWater", 128, 128, 64, 0, 1.0).Appends(
		floortemplate.SoilWater128x128()...,
	)
	tw.Add("SoilMagma", 128, 128, 64, 0, 1.0).Appends(
		floortemplate.SoilMagma128x128()...,
	)
	tw.Add("SoilIce", 128, 128, 64, 0, 1.0).Appends(
		floortemplate.SoilIce128x128()...,
	)
	tw.Add("MadeByImage", 256, 256, 256, 0, 1.0).Appends(
		"ResourceFromPNG name=imagefloor.png",
		"ResourceRand resource=Plant mean=100000000 stddev=65535 repeat=256",
		"ResourceAgeing initrun=0 msper=64000 resetaftern=1440",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=32 mean=6 stddev=4 min=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	)
	tw.Add("AgeingCity", 256, 256, 256, 0, 1.0).Appends(
		floortemplate.AgeingCity256x256()...,
	)
	tw.Add("AgeingField", 256, 256, 256, 0, 1.0).Appends(
		floortemplate.AgeingField256x256()...,
	)
	tw.Add("AgeingMaze", 256, 256, 256, 0, 1.0).Appends(
		floortemplate.AgeingMaze256x256()...,
	)

	fm := tw.Add("BedTown", 256, 256, 128, 0, 1.0).Appends(
		floortemplate.CityRooms(256, 256, 11, 11, 5, rnd.Intn)...,
	)

	fm = tw.Add("ResourceMazeFill", 256, 256, 128, 0, 1.0).Appends(
		floortemplate.MixedResourceMaze(256, 256)...,
	)
	fm.Appends(
		fmt.Sprintf("ResourceFillRect resource=Soil  amount=1  x=0 y=0  w=%v h=%v", 256, 256),
	)
	fm.Appends(
		floortemplate.CityRoomsRand(512, rnd.Intn)...,
	)

	fm = tw.Add("ResourceMaze", 256, 256, 128, 0, 1.0).Appends(
		floortemplate.MixedResourceMaze(256, 256)...,
	)
	fm.Appends(
		floortemplate.CityRoomsRand(512, rnd.Intn)...,
	)

	tw.Add("RogueLike", 80, 43, 16, 0, 1.0).Appends(
		floortemplate.RogueLike80x43()...,
	)
	tw.Add("GogueLike", 80, 43, 32, 0, 1.0).Appends(
		floortemplate.GogueLike()...,
	)
	tw.Add("Ghost", 80, 43, 16, 0, 1.0).Appends(
		floortemplate.Ghost80x43()...,
	)
	tw.Add("FreeForAll", 64, 64, 16, 0, 1.0).Appends(
		floortemplate.FreeForAll64x64()...,
	)
	tw.Add("TileRooms", 64, 32, 0, 0, 1.5).Appends(
		floortemplate.TileRooms64x32()...,
	)
	tw.Add("PortalMaze", 64, 32, 0, 0, 1.5).Appends(
		floortemplate.PortalMaze64x32Finalized()...,
	)
	tw.Add("MazeRooms1", 64, 32, 0, 0, 1.5).Appends(
		floortemplate.MazeBigSmall64x32()...,
	)
	tw.Add("MazeRooms2", 64, 32, 0, 0, 1.5).Appends(
		floortemplate.MazeRooms64x32()...,
	)
	tw.Add("MazeRooms3", 64, 32, 0, 0, 1.5).Appends(
		floortemplate.MazeRoomsOverlapWall64x32()...,
	)
	tw.Add("MazeWalk", 64, 64, 0, 0, 2.0).Appends(
		"ResourceMazeWalk resource=Soil amount=64 x=0 y=0 w=64 h=64 xn=16 yn=16 connerfill=true",
	)

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}

	tw.GetByName("ManyPortals").ConnectPortalTo("Rand", "Rand", tw.GetByName("TileRooms"))
	tw.GetByName("TileRooms").ConnectPortalTo("Rand", " x=07 y=07", tw.GetByName("PortalMaze"))
	tw.GetByName("PortalMaze").ConnectPortalTo(" x=55 y=23", " x=15 y=15", tw.GetByName("MazeRooms1"))
	tw.GetByName("MazeRooms1").ConnectPortalTo(" x=47 y=15", "InRoom", tw.GetByName("MazeRooms2"))
	tw.GetByName("MazeRooms2").ConnectPortalTo("Rand", "Rand", tw.GetByName("MazeRooms3"))
	tw.GetByName("MazeRooms3").ConnectPortalTo("Rand", "Rand", tw.GetByName("MazeWalk"))
	tw.GetByName("MazeWalk").ConnectPortalTo("Rand", "Rand", tw.GetByName("ManyPortals"))

	tw.GetByName("Practice").ConnectStairUp("InRoom", "InRoom", tw.GetByName("SoilPlant"))
	tw.GetByName("SoilPlant").ConnectStairUp("InRoom", "InRoom", tw.GetByName("AgeingCity"))
	tw.GetByName("AgeingCity").ConnectStairUp("InRoom", "InRoom", tw.GetByName("BedTown"))
	tw.GetByName("BedTown").ConnectStairUp("InRoom", "InRoom", tw.GetByName("RogueLike"))

	tw.GetByName("RogueLike").ConnectStairUp("InRoom", "InRoom", tw.GetByName("SoilWater"))
	tw.GetByName("SoilWater").ConnectStairUp("InRoom", "InRoom", tw.GetByName("AgeingField"))
	tw.GetByName("AgeingField").ConnectStairUp("InRoom", "InRoom", tw.GetByName("ResourceMazeFill"))
	tw.GetByName("ResourceMazeFill").ConnectStairUp("InRoom", "InRoom", tw.GetByName("GogueLike"))

	tw.GetByName("GogueLike").ConnectStairUp("InRoom", "InRoom", tw.GetByName("SoilMagma"))
	tw.GetByName("SoilMagma").ConnectStairUp("InRoom", "InRoom", tw.GetByName("AgeingMaze"))
	tw.GetByName("AgeingMaze").ConnectStairUp("InRoom", "InRoom", tw.GetByName("ResourceMaze"))
	tw.GetByName("ResourceMaze").ConnectStairUp("InRoom", "InRoom", tw.GetByName("Ghost"))

	tw.GetByName("Ghost").ConnectStairUp("InRoom", "Rand", tw.GetByName("SoilIce"))
	tw.GetByName("SoilIce").ConnectStairUp("InRoom", "InRoom", tw.GetByName("MadeByImage"))
	tw.GetByName("MadeByImage").ConnectStairUp("InRoom", "Rand", tw.GetByName("FreeForAll"))

	tw.GetByName("FreeForAll").ConnectStairUp("Rand", "InRoom", tw.GetByName("Practice"))

	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("Practice"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilPlant"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilWater"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilMagma"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilIce"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("MadeByImage"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("AgeingCity"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("AgeingField"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("AgeingMaze"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("RogueLike"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("GogueLike"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("Ghost"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("FreeForAll"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("BedTown"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("ResourceMazeFill"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("ResourceMaze"))

	for _, fm := range tw.GetList() {
		randList := map[string]bool{
			"FreeForAll": true,
			"PortalMaze": true,
			"MazeRooms1": true,
			"MazeRooms2": true,
			"MazeRooms3": true,
			"MazeWalk":   true,
		}
		suffix := "InRoom"
		if randList[fm.Name] {
			suffix = "Rand"
		}
		fm.ConnectAutoInPortalTo(suffix, suffix, fm)
		fm.AddTrapTeleportTo(suffix, fm)
		fm.AddAllEffectTrap(suffix, 1)

		roomCount := fm.CalcRoomCount()
		// fmt.Printf("%v Room %v\n", fm, roomCount)
		recycleCount := fm.W * fm.H / 512
		if recycleCount < 2 {
			recycleCount = 2
		}
		if recycleCount > roomCount {
			fm.AddRecycler(suffix, roomCount)
			fm.AddRecycler("Rand", recycleCount-roomCount)
		} else {
			fm.AddRecycler(suffix, recycleCount)
		}
	}

	return tw
}
