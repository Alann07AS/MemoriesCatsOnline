package randomcat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetXRandomCatImg(nb int) (table []string) {
	for nb != 0 {
		time.Sleep(50 * time.Millisecond)
		r, err := http.Get("https://apilist.fun/out/randomcat")
		if err != nil {
			fmt.Println("err:", err)
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("err:", err)
		}
		f := Fille{}
		err = json.Unmarshal(content, &f)
		if err != nil {
			fmt.Println("err:", err)
		}
		table = append(table, f.File)
		nb--
	}
	fmt.Println("GetCardFromAPIDone")
	return table
}

type Fille struct {
	File string
}
