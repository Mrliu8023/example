package segmenter

import (
	groups_dev "example/groups-dev"
	"example/search"
	"fmt"
	"github.com/go-ego/gpy"
	"github.com/go-ego/gse"
	"strings"
	"testing"
)

func TestGseSegment(t *testing.T) {
	gset, _ := gse.New("./dict.txt")
	text := "基础设施/综合布线/其他/维谛/维谛/公共模块/温湿度th03"

	//fmt.Println("gse ********************")
	//fmt.Println("全模式： ", gset.CutAll(text))
	//fmt.Println("基础模式： ", gset.Cut(text))
	//fmt.Println("搜索模式: ", gset.CutSearch(text))
	//segs := gset.ModeSegment([]byte(text), false)
	//for _, s := range segs {
	//	fmt.Println(s.Token().Text())
	//}
	fmt.Println(PinYin(text, &gset))

}

// PinYin get the Chinese alphabet and abbreviation
func PinYin(hans string, sgt *gse.Segmenter) []string {

	var (
		str      string
		pyStr    string
		strArr   []string
		pyArr    []string
		splitStr string
		// splitArr []string
	)

	//
	splitHans := strings.Split(hans, "")
	for i := 0; i < len(splitHans); i++ {
		if splitHans[i] != "" {
			strArr = append(strArr, splitHans[i])
			splitStr += splitHans[i]
			strArr = append(strArr, splitStr)
		}
	}

	// Segment 分词

	sehans := sgt.CutSearch(hans, true)
	for h := 0; h < len(sehans); h++ {
		strArr = append(strArr, sehans[h])
	}

	//
	// py := pinyin.LazyConvert(sehans[h], nil)
	pyMap := make(map[string]struct{})
	// fmt.Println(strArr)
	py := gpy.LazyConvert(hans, nil)

	// fmt.Println("py...", py)
	for i := 0; i < len(py); i++ {
		// log.Println("py[i]...", py[i])
		pyStr += py[i]

		pyMap[pyStr] = struct{}{}
		pyArr = append(pyArr, pyStr)

		if len(py[i]) > 0 {
			str += py[i][0:1]

			pyMap[pyStr] = struct{}{}
			pyArr = append(pyArr, str)

		}
	}

	for _, han := range strArr {
		str = ""
		py = gpy.LazyConvert(han, nil)
		// fmt.Println("py: ", py)
		for i := 0; i < len(py); i++ {
			if _, ok := pyMap[py[i]]; !ok {
				pyMap[py[i]] = struct{}{}
				pyArr = append(pyArr, py[i])
			}
			if len(py[i]) > 0 {
				str += py[i][0:1]

				if _, ok := pyMap[str]; !ok {
					pyMap[py[i]] = struct{}{}
					pyArr = append(pyArr, str)
				}

			}
		}
		// fmt.Println("pyArr: ", pyArr)
	}
	strArr = append(strArr, pyArr...)

	return strArr
}

func TestNewSegmenter(t *testing.T) {
	seg, err := NewSegmenter(true, true, true)
	if err != nil {
		t.Fatal(err)
	}

	gl, err := groups_dev.Parse("D:\\go\\src\\example\\groups_dev.json")
	if err != nil {
		t.Fatal(err)
	}

	docs := make([]search.Document, 0, len(gl.Groups))
	for i, g := range gl.Groups {
		if i < 10 {
			docs = append(docs, &Document{
				id:      g.ID,
				content: g.Display,
			})
		}
	}

	seg.AddDocuments(docs)

	//gl, err = seg.Search("wsd")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//for _, g := range gl.Groups {
	//	fmt.Println(g)
	//}

}
