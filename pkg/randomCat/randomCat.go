package randomcat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const APIurl = "https://api.thecatapi.com/v1/images/search"

func GetXRandomCatImg(nb int) (table []string) {
	for nb != 0 {
		time.Sleep(50 * time.Millisecond)
		r, err := http.Get(APIurl)
		if err != nil {
			fmt.Println("err:", err)
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("err:", err)
		}
		f := make([]Fille, 1)
		err = json.Unmarshal(content, &f)
		if err != nil {
			fmt.Println("err:", err)
		}
		table = append(table, f[0].Url)
		nb--
	}
	fmt.Println("GetCardFromAPIDone")
	return table
}

type Fille struct {
	Url string `json:"url"`
}
