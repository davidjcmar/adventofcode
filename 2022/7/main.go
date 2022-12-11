package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type file_p struct {
	size int
	name string
}

type dir_p struct {
	name   string
	parent *dir_p
	dirs   map[string]*dir_p
	files  map[string]*file_p
}

type root struct {
	root *dir_p
}

func newFilesystem() *dir_p {
	f := dir_p{name: "/", parent: nil, dirs: nil, files: nil}
	return &f
}

func newDir(name string, parent *dir_p) *dir_p {
	d := dir_p{name: name, parent: parent}
	return &d
}

func newFile(name string, size int) *file_p {
	f := file_p{name: name, size: size}
	return &f
}

func buildFilesystem(scan *bufio.Scanner, root *dir_p) *dir_p {
	current_dir := root
	for scan.Scan() {
		line := scan.Text()
		//fmt.Printf("Runline: %s\n", line)

		if line[0] != '$' {
			ls(current_dir, root, line)
		} else if strings.HasPrefix(line, "$") && strings.Contains(line, "cd") {
			//fmt.Printf("In buildfs c_dir: %s dir: %s\n", current_dir.name, line[5:])
			current_dir = cd(current_dir, root, line)
		}
	}
	return root
}

func print(filename, filetype string, size, depth int) {
	var ws string
	for i := 0; i < depth; i++ {
		ws += "  "
	}
	if filetype == "dir" {
		fmt.Printf("%s- %s (%s)\n", ws, filename, filetype)
	} else {
		fmt.Printf("%s- %s (%s, size=%d)\n", ws, filename, filetype, size)
	}
}

func printFilesystem(d *dir_p, depth int) {
	var ws string
	for i := 0; i < depth; i++ {
		ws += "  "
	}
	print(d.name, "dir", 0, depth)
	for _, dir := range d.dirs {
		printFilesystem(dir, depth+1)
	}
	for _, file := range d.files {
		print(file.name, "file", file.size, depth+1)
	}
}

func get_dirs(c_dir *dir_p) []string {
	var dirs []string
	for _, d := range c_dir.dirs {
		dirs = append(dirs, d.name)
	}
	return dirs
}
func cd(c_dir *dir_p, root *dir_p, command string) *dir_p {
	prefix := "$ cd "
	t_dir := strings.TrimPrefix(command, prefix)
	switch t_dir {
	case "/":
		//fmt.Printf("in cd return root\n")
		return root
	case "..":
		//fmt.Printf("in cd return parent: %s\n", c_dir.parent.name)
		return c_dir.parent
	default:
		if c_dir.dirs == nil {
			//fmt.Printf("in cd init dirs c_dir: %s\n", c_dir.name)
			c_dir.dirs = make(map[string]*dir_p)
		}
		if dir, ok := c_dir.dirs[t_dir]; ok {
			//fmt.Printf("in cd existing dir c_dir: %s t_dir: %s\n", c_dir.name, t_dir)
			return dir
		} else {
			//fmt.Printf("in cd new dir c_cir: %s t_dir: %s\n", c_dir.name, t_dir)
			c_dir.dirs[t_dir] = newDir(t_dir, c_dir)
			//fmt.Printf("c_dir: %s -- new dirs: %v\n", c_dir.name, c_dir.dirs)
			return c_dir.dirs[t_dir]
		}
	}
}

func ls(c_dir, root *dir_p, line string) {
	switch line[:3] {
	case "dir":
		// directory
		//fmt.Printf("in ls dir: %s c_dir: %s t_dir: %s\n", line, c_dir.name, line[4:])
		c_dir.dirs[line[4:]] = cd(c_dir, root, line[4:])
	default:
		// file
		//fmt.Printf("in ls file: %s c_dir: %s\n", line, c_dir.name)
		fields := strings.Fields(line)
		name := fields[1]
		size, _ := strconv.Atoi(fields[0])
		if c_dir.files == nil {
			c_dir.files = make(map[string]*file_p)
		}
		c_dir.files[name] = newFile(name, size)
	}
}

func sumFirstPart(root *dir_p) int {
	fmt.Println(root)
	return 0
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bufio.NewReader(f))
	root := newFilesystem()
	root = buildFilesystem(scan, root)
	f.Close()
	printFilesystem(root, 0)
	//fmt.Printf("root dirs: %v\n", root.dirs)
}
