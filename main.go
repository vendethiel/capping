package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Character struct {
	Name    string    `json:"name"`
	Class   string    `json:"class"`
	LastCap time.Time `json:"last_cap"`
}

type data struct {
	Characters []Character `json:"characters"`
}

func isCapped(character *Character) bool {
	t := time.Now()
	year, month, day := t.Date()
	lastWed := time.Date(year, month, day, 3, 0, 0, 0, t.Location())
	for lastWed.Weekday() != time.Wednesday {
		lastWed = lastWed.AddDate(0, 0, -1)
	}
	return character.LastCap.After(lastWed)
}

func capped(character *Character) string {
	if isCapped(character) {
		return "capped"
	} else {
		return "NOT CAPPED"
	}
}

func normalizeClass(c string) string {
	switch c {
	case "rog", "roug", "rouge":
		return "rogue"
	case "dh", "demon", "dhunter", "demonh":
		return "demonhunter"
	case "dk", "deathk", "dknight", "deathnigt", "deathniht":
		return "deathknight"
	case "pal", "pally", "paly", "paladino":
		return "paladin"
	case "war", "warrior", "warr", "krieger":
		return "warrior"
	case "dru", "drui", "druide":
		return "druid"
	case "mnk":
		return "monk"
	case "hunt":
		return "hunter"
	case "sham", "cham", "chaman":
		return "shaman"
	case "mag", "maeg":
		return "mage"
	case "pri", "prie", "prist":
		return "priest"
	case "dragon", "evo", "evokr":
		return "evoker"
	case "wlock", "lock", "lck":
		return "warlock"
	default:
		return c
	}
}

func compareClass(a, b string) bool {
	return strings.EqualFold(normalizeClass(a), normalizeClass(b))
}

func match(character *Character, find string) bool {
	return strings.EqualFold(character.Name, find) || compareClass(character.Class, find)
}

func locate(d *data, find string) (int, string) {
	idx := -1
	for i := 0; i < len(d.Characters); i++ {
		if match(&d.Characters[i], find) {
			if idx == -1 {
				idx = i
			} else {
				return -1, "Multiple characters matched"
			}
		}
	}
	if idx == -1 {
		return -1, "No such character"
	}
	return idx, ""
}

func main() {
	f, err := os.ReadFile("chars.json")
	if err != nil {
		panic("Couldn't open file")
	}

	var d data
	if err := json.Unmarshal(f, &d); err != nil {
		panic("Couldn't parse file")
	}

	if len(os.Args) < 2 {
		hasLeft := false
		for i := 0; i < len(d.Characters); i++ {
			if !isCapped(&d.Characters[i]) {
				if !hasLeft {
					hasLeft = true
					println("Characters left to cap:")
				}

				fmt.Printf("- %s: %s\n", d.Characters[i].Name, d.Characters[i].Class)
			}
		}
		if !hasLeft {
			println("All characters capped :)")
		}
	} else {
		switch os.Args[1] {

		case "cap":
			if len(os.Args) != 3 {
				panic("You need to provide a character to cap")
			}

			find := os.Args[2]
			idx, err := locate(&d, find)
			if err != "" {
				panic(err)
			}
			d.Characters[idx].LastCap = time.Now()

		case "add":
			if len(os.Args) != 4 {
				panic("You need to provide charname + name")
			}

			name := os.Args[2]
			class := normalizeClass(os.Args[3])
			for i := 0; i < len(d.Characters); i++ {
				if strings.EqualFold(d.Characters[i].Name, name) {
					panic("A character with that name already exists")
				}
			}
			d.Characters = append(d.Characters, Character{Name: name, Class: class, LastCap: time.Now()})

		case "list":
			println("Characters:")
			for i := 0; i < len(d.Characters); i++ {
				fmt.Printf("- %s: %s (%s)\n", d.Characters[i].Name, d.Characters[i].Class, capped(&d.Characters[i]))
			}

		default:
			find := os.Args[1]
			idx, err := locate(&d, find)
			if err != "" {
				panic("Unknown command/character")
			}
			d.Characters[idx].LastCap = time.Now()
			fmt.Printf("Capped %s\n", d.Characters[idx].Name)
		}
	}

	w, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		panic("Couldn't encode json")
	}

	if err := ioutil.WriteFile("chars.json", w, 0); err != nil {
		panic("Couldn't write the file")
	}
}
