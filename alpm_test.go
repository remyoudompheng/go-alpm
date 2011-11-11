package main

import (
	. "alpm"
	"fmt"
	"os"
)

const (
	root    = "/"
	dbpath  = "/var/lib/pacman"
	version = "7.0.0"
)

var h *Handle

func init() {
	var err os.Error
	h, err = Init("/", "/var/lib/pacman")
	if err != nil {
		fmt.Printf("failed to Init(): %s", err)
		os.Exit(1)
	}
}

func TestVersion() {
	if Version() != version {
		fmt.Println("verion's do not match")
	}
}

func TestVercmp() {
	x := VerCmp("1.0-2", "2.0-1")
	if x >= 0 {
		fmt.Println("failed at checking 2.0-1 is newer than 1.0-2")
	}
	x = VerCmp("1:1.0-2", "2.0-1")
	if x <= 0 {
		fmt.Println("failed at checking 2.0-1 is older than 1.0-2")
	}
	x = VerCmp("2.0.2-2", "2.0.2-2")
	if x != 0 {
		fmt.Println("failed at checking 2.0.2-2 is equal to itself")
	}
}

func TestRevdeps() {
	fmt.Print("Testing reverse deps of glibc...\n")
	db, _ := h.LocalDb()
	pkg, _ := db.GetPkg("glibc")
	for _, pkgname := range pkg.ComputeRequiredBy() {
		fmt.Println(pkgname)
	}
}

func TestLocalDB() {
	defer func() {
		if recover() != nil {
			fmt.Println("local db failed")
		}
	}()
	db, _ := h.LocalDb()
	fmt.Print("Testing listing local db...\n")
	number := 0
	for pkg := range db.PkgCache() {
		number++
		if number <= 15 {
			fmt.Printf("%v \n", pkg.Name())
		}
	}
	if number > 15 {
		fmt.Printf("%d more packages...\n", number-15)
	}
}

func TestRelease() {
	if err := h.Release(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	TestVersion()
	TestVercmp()
	TestRevdeps()
	TestLocalDB()
	TestRelease()
}
