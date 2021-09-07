package groups_dev

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Group struct {
	Display string `json:"display"`
	ID      string `json:"id"`
	Value   string `json:"value"`
}

type GroupList struct {
	Groups []*Group
	Length int
}

func Parse(fileName string) (gl *GroupList, err error) {
	var gPathMap = map[string]string{}
	var gMap = map[string]string{}

	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var gs []*Group

	if err := json.Unmarshal(bs, &gs); err != nil {
		return nil, err
	}

	var p string
	for _, g := range gs {
		if strings.LastIndex(g.ID, ".") == -1 {
			p = g.Value
		} else {
			v, ok := gPathMap[g.ID[:strings.LastIndex(g.ID, ".")]]
			if ok {
				p = v + "/" + g.Value
			} else {
				var gp []string
				for i := 0; i < len(g.ID)-1; i++ {
					if rune(g.ID[i]) == '.' {
						v := gMap[g.ID[:i]]
						gp = append(gp, v)
					}
				}
				gp = append(gp, g.Value)
				gPathMap[g.ID] = strings.Join(gp, "/")

				p = strings.Join(gp, "/")
			}
		}
		g.Display = p
		gPathMap[g.ID] = p
	}

	gl = &GroupList{
		Groups: gs,
		Length: len(gs),
	}

	return gl, nil
}
