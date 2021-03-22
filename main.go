package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	humanize "github.com/dustin/go-humanize"
	godirwalk "github.com/karrick/godirwalk"
	str2duration "github.com/xhit/go-str2duration/v2"
)

var inodeset = make(map[uint64]struct{})
var sizestats = make(map[bool]uint64)

func main() {
	path := flag.String("path", ".", "path where to count for unique inodes")
	olderthan := flag.String("olderthan", "365d", "count file size older/more recent than the specified interval")
	flag.Parse()
	abspath, err := filepath.Abs(*path)

	interval, err := str2duration.ParseDuration(*olderthan)
	if err != nil {
		panic(err)
	}

	err = godirwalk.Walk(abspath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				stats, errstat := os.Stat(osPathname)
				st := stats.Sys().(*syscall.Stat_t)
				if errstat != nil {
					fmt.Println("stat failed:", errstat)
				} else {
					if _, ok := inodeset[st.Ino]; ! ok {
						inodeset[st.Ino] = struct{}{}
						mtime := stats.ModTime()
						isitold := time.Now().Sub(mtime) > interval
						sizestats[isitold] += uint64(stats.Size())
					}
				}
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	if err != nil {
		log.Fatal("walk failed: ", err)
	}

	fmt.Printf("%s: %d unique inodes\n", abspath, len(inodeset))
	fmt.Printf("%s: > %s: %s < %s: %s\n", abspath, *olderthan, humanize.Bytes(sizestats[false]), *olderthan, humanize.Bytes(sizestats[true]))
}
