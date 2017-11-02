package main

import (
	"github.com/google/flatbuffers/go"
	"github.com/whitecreepy/flatbuffer-test/MyGame/Sample"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"log"
	"os"
)

func generateMonster() []byte {
	builder := flatbuffers.NewBuilder(1024)

	weaponOne := builder.CreateString("Sword")
	weaponTwo := builder.CreateString("Axe")

	// Create sword
	Sample.WeaponStart(builder)
	Sample.WeaponAddName(builder, weaponOne)
	Sample.WeaponAddDamage(builder, 3)
	sword := Sample.WeaponEnd(builder)

	// Create axe
	Sample.WeaponStart(builder)
	Sample.WeaponAddName(builder, weaponTwo)
	Sample.WeaponAddDamage(builder, 5)
	axe := Sample.WeaponEnd(builder)

	// Serialize a name for our monster, called "Orc".
	name := builder.CreateString("Orc")
	// Create a `vector` representing the inventory of the Orc. Each number
	// could correspond to an item that can be claimed after he is slain.
	// Note: Since we prepend the bytes, this loop iterates in reverse.
	Sample.MonsterStartInventoryVector(builder, 10)
	for i := 9; i >= 0; i-- {
		builder.PrependByte(byte(i))
	}
	inv := builder.EndVector(10)

	// Create a FlatBuffer vector and prepend the weapons.
	// Note: Since we prepend the data, prepend them in reverse order.
	Sample.MonsterStartWeaponsVector(builder, 2)
	builder.PrependUOffsetT(axe)
	builder.PrependUOffsetT(sword)
	weapons := builder.EndVector(2)

	Sample.MonsterStartPathVector(builder, 2)
	Sample.CreateVec3(builder, 1.0, 2.0, 3.0)
	Sample.CreateVec3(builder, 4.0, 5.0, 6.0)
	path := builder.EndVector(2)

	// Create our monster using `MonsterStart()` and `MonsterEnd()`.
	Sample.MonsterStart(builder)
	Sample.MonsterAddPos(builder, Sample.CreateVec3(builder, 1.0, 2.0, 3.0))
	Sample.MonsterAddHp(builder, 300)
	Sample.MonsterAddName(builder, name)
	Sample.MonsterAddInventory(builder, inv)
	Sample.MonsterAddColor(builder, Sample.ColorRed)
	Sample.MonsterAddWeapons(builder, weapons)
	Sample.MonsterAddEquippedType(builder, Sample.EquipmentWeapon)
	Sample.MonsterAddEquipped(builder, axe)
	Sample.MonsterAddPath(builder, path)
	orc := Sample.MonsterEnd(builder)


	builder.Finish(orc)

	buf := builder.FinishedBytes()
	return buf
}

func writeMonsterToFile(path string, monster []byte) {
	err := ioutil.WriteFile("orc.txt", monster, 0644)
	if err != nil {
		panic(err)
	}
}

func readMonsterFromFile(path string) *Sample.Monster {
	buf, err := ioutil.ReadFile("empty.txt")

	if err != nil {
		panic(err)
	}

	monster := Sample.GetRootAsMonster(buf, 0)
	return monster
}

func monsterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(generateMonster())
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api", monsterHandler).Methods("GET", "POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("simpleUI/static")))

	err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router))

	log.Println("Listening at port 8080.")
	if err != nil {
		log.Fatal("Listen and serve failed: ", err)
	}
}
