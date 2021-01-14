package main

import (
	"fmt"
	"flag"
	"log"
	"syscall"
	"path/filepath"

	"github.com/karrick/godirwalk"
)

var inodeset = make(map[uint64]struct{})


func main() {
	path := flag.String("path", ".", "path where to count for unique inodes")
	flag.Parse()
	abspath, err := filepath.Abs(*path)
	
	err = godirwalk.Walk(abspath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				var stat syscall.Stat_t
				//fmt.Println("osPath:", osPathname)
				errstat := syscall.Stat(osPathname, &stat)
				if errstat != nil {
					fmt.Println("stat failed:", errstat)
				} else {
					inodeset[stat.Ino] = struct{}{}					
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
}