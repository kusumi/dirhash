package main

import (
	"fmt"
)

var (
	statDirectory   []string // hashed
	statRegular     []string // hashed
	statDevice      []string // hashed
	statSymlink     []string // hashed
	statUnsupported []string
	statInvalid     []string
	statIgnored     []string

	writtenDirectory uint // hashed
	writtenRegular   uint // hashed
	writtenDevice    uint // hashed
	writtenSymlink   uint // hashed
)

func initStat() {
	statDirectory = make([]string, 0)
	statRegular = make([]string, 0)
	statDevice = make([]string, 0)
	statSymlink = make([]string, 0)
	statUnsupported = make([]string, 0)
	statInvalid = make([]string, 0)
	statIgnored = make([]string, 0)

	writtenDirectory = 0
	writtenRegular = 0
	writtenDevice = 0
	writtenSymlink = 0
}

// num stat
func numStatTotal() uint {
	return numStatDirectory() + numStatRegular() + numStatDevice() + numStatSymlink()
}

func numStatDirectory() uint {
	return uint(len(statDirectory))
}

func numStatRegular() uint {
	return uint(len(statRegular))
}

func numStatDevice() uint {
	return uint(len(statDevice))
}

func numStatSymlink() uint {
	return uint(len(statSymlink))
}

/*
func numStatUnsupported() uint {
	return uint(len(statUnsupported))
}

func numStatInvalid() uint {
	return uint(len(statInvalid))
}

func numStatIgnored() uint {
	return uint(len(statIgnored))
}
*/

// append stat
func appendStatTotal() {
}

func appendStatDirectory(f string) {
	statDirectory = append(statDirectory, f)
}

func appendStatRegular(f string) {
	statRegular = append(statRegular, f)
}

func appendStatDevice(f string) {
	statDevice = append(statDevice, f)
}

func appendStatSymlink(f string) {
	statSymlink = append(statSymlink, f)
}

func appendStatUnsupported(f string) {
	statUnsupported = append(statUnsupported, f)
}

func appendStatInvalid(f string) {
	statInvalid = append(statInvalid, f)
}

func appendStatIgnored(f string) {
	statIgnored = append(statIgnored, f)
}

// print stat
/*
func printStatDirectory() {
	printStat(statDirectory, strDir)
}

func printStatRegular() {
	printStat(statRegular, strReg)
}

func printStatDevice() {
	printStat(statDevice, strDevice)
}

func printStatSymlink() {
	printStat(statSymlink, strSymlink)
}
*/

func printStatUnsupported() {
	printStat(statUnsupported, strUnsupported)
}

func printStatInvalid() {
	printStat(statInvalid, strInvalid)
}

func printStatIgnored() {
	printStat(statIgnored, "ignored file")
}

func printStat(l []string, msg string) {
	if len(l) == 0 {
		return
	}
	printNumFormatString(uint(len(l)), msg)

	for _, v := range l {
		f := getRealPath(v)
		t1, _ := getRawFileType(v)
		t2, _ := getFileType(v)
		assert(t2 != typeSymlink) // symlink chains resolved
		if t1 == typeSymlink {
			assert(optIgnoreSymlink || t2 == typeDir || t2 == typeInvalid)
			fmt.Printf("%s (%s -> %s)\n",
				f, getFileTypeString(t1), getFileTypeString(t2))
		} else {
			assert(t2 != typeDir)
			fmt.Printf("%s (%s)\n", f, getFileTypeString(t1))
		}
	}
}

// num written
func numWrittenTotal() uint {
	return numWrittenDirectory() + numWrittenRegular() + numWrittenDevice() + numWrittenSymlink()
}

func numWrittenDirectory() uint {
	return writtenDirectory
}

func numWrittenRegular() uint {
	return writtenRegular
}

func numWrittenDevice() uint {
	return writtenDevice
}

func numWrittenSymlink() uint {
	return writtenSymlink
}

// append written
func appendWrittenTotal(written uint64) {
}

func appendWrittenDirectory(written uint64) {
	writtenDirectory += uint(written)
}

func appendWrittenRegular(written uint64) {
	writtenRegular += uint(written)
}

func appendWrittenDevice(written uint64) {
	writtenDevice += uint(written)
}

func appendWrittenSymlink(written uint64) {
	writtenSymlink += uint(written)
}
