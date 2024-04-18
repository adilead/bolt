package bookmarks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/adilead/bolt/internal/bolt"
)


type BookmarkItem struct {
    name string 
    url string
}

func (bi *BookmarkItem) ToSlice() []string {
    return []string{bi.name, bi.url}
}

type BookmarkRoot struct {
    Checksum string `json:"checksum"`
    Roots map[string]BookmarkCat `json:"roots"`
}

type BookmarkCat struct {
    Children []Bookmark `json:"children"`
}

type Bookmark struct {
    Date_added string `json:"date_added"`
    Date_last_added string `json:"date_last_added"`
    Guid string `json:"guid"`
    Id string `json:"id"`
    Name string `json:"name"`
    BookmarkType string `json:"type"`
    Url string `json:"url"`
}

type Bookmarks struct {
    
}

func GetBookmarks () ([]bolt.Searchable, error) {
    bookmark_path := "/home/adrian/.config/BraveSoftware/Brave-Browser/Default/Bookmarks"
    f, err := os.Open(bookmark_path)
    if err != nil {
        return nil, err
    }
        
    // var data map[string]interface{}
    decoder := json.NewDecoder(bufio.NewReader(f))
    br := BookmarkRoot{}
    err = decoder.Decode(&br)
    if err != nil {
        return nil, err
    }
    fmt.Println(br.Roots["bookmark_bar"].Children[0])


    bookmark_items := make([]bolt.Searchable, 0)

    for _, bookmark := range br.Roots["bookmark_bar"].Children {
        name := bookmark.Name
        url := bookmark.Url
        if len(name) == 0 {
            name = url
        }         
        fmt.Println(name)
        bookmark_items = append(bookmark_items, &BookmarkItem{name: name, url: url})
    }

    return bookmark_items, nil
}



func OpenBookmark(choice bolt.Searchable) error {

    bi := choice.(*BookmarkItem)
    cmd := exec.Command("brave", "--new-window", bi.name)
    cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
    err := cmd.Start()
    if err != nil {
        panic(err)
    }
    return nil
} 
