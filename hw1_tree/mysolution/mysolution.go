package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

// ReadDir ..
func ReadDir(b *bytes.Buffer, dirname string, verbose bool, dirPrfx string) (bytes.Buffer, error) {
	res, vrb := b, verbose

	corner, taur, mystr := "└───", "├───", ""
	dirCount, fileCount := 0, 0
	newDirPrfx, fileSize := "", ""

	f, err := os.Open(dirname)
	if err != nil {
		return *res, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return *res, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	// Counting dirs/files in filelist in order to connect vertical line
	// from top to latest dir, but not lower.
	for _, f := range list {
		if f.IsDir() {
			dirCount++
		} else if verbose == true {
			fileCount++
		}
	}

	for _, f := range list {
		if dirCount == 1 && fileCount == 0 {
			newDirPrfx = dirPrfx + "\t"
		} else {
			newDirPrfx = dirPrfx + "│\t"
		}

		switch f.IsDir() {
		case true:
			if dirCount > 1 {
				mystr = fmt.Sprintf("%s%s", dirPrfx, taur)
				dirCount--
			} else if dirCount == 1 && fileCount > 0 {
				mystr = fmt.Sprintf("%s%s", dirPrfx, taur)
				dirCount--
			} else {
				mystr = fmt.Sprintf("%s%s", dirPrfx, corner)
			}
			res.WriteString(fmt.Sprintf("%s%s\n", mystr, f.Name()))
			*res, err = ReadDir(res, dirname+"/"+f.Name(), vrb, newDirPrfx)
		case false:
			if f.Size() > 0 {
				fileSize = fmt.Sprintf("%db", f.Size())
			} else {
				fileSize = "empty"
			}
			if fileCount == 1 && dirCount == 0 {
				mystr = fmt.Sprintf("%s%s%s (%s)\n", dirPrfx, corner, f.Name(), fileSize)
				fileCount--
			} else if fileCount == 1 && dirCount > 0 {
				mystr = fmt.Sprintf("%s%s%s (%s)\n", dirPrfx, taur, f.Name(), fileSize)
				fileCount--
			} else if fileCount > 1 {
				mystr = fmt.Sprintf("%s%s%s (%s)\n", dirPrfx, taur, f.Name(), fileSize)
				fileCount--
			}
			if verbose == true {
				res.WriteString(mystr)
			}
		}
	}
	return *res, nil
}

func dirTree(out *bytes.Buffer, data string, verbose bool) error {
	_, err := ReadDir(out, data, verbose, "")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	buf := new(bytes.Buffer)
	args := os.Args
	var a bytes.Buffer
	var err error

	if len(args) == 1 {
		a, err = ReadDir(buf, ".", false, "")
	} else if len(args) == 2 {
		a, err = ReadDir(buf, args[1], false, "")
	} else if len(args) > 2 && args[2] == "-f" {
		a, err = ReadDir(buf, args[1], true, "")
	} else {
		fmt.Printf("Invalid format.\nformat: main.go [path] [-f]")
	}

	if err != nil {
		fmt.Println("Something went wrong")
	}

	fmt.Print(a.String())
}
