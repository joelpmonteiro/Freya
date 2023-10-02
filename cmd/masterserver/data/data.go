package data

import (
	"io/ioutil"

	"github.com/ubis/Freya/share/directory"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/models/skills"

	"gopkg.in/yaml.v2"
)

type Loader struct {
	*InitialData
}

type InitialData struct {
	BattleStyles []struct {
		ID        int
		Location  map[string]int
		Stats     map[string]int
		Equipment map[string]inventory.Item
		Inventory map[int]inventory.Item
		Skills    map[int]skills.Skill
		Links     map[int]skills.Link
	}
}

// Initializes DataLoader
func (dl *Loader) Init() {
	log.Info("Loading data...")

	dl.InitialData = &InitialData{}
	dl.load("initial_data.yml", dl.InitialData)
}

// Deserializes data from file to specified struct
func (dl *Loader) load(filename string, data interface{}) {
	var s, err = ioutil.ReadFile(directory.Root() + "/data/" + filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = yaml.Unmarshal(s, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}
