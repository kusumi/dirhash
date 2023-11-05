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

	writtenDirectory uint64 // hashed
	writtenRegular   uint64 // hashed
	writtenDevice    uint64 // hashed
	writtenSymlink   uint64 // hashed
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
func numStatTotal() uint64 {
	return numStatDirectory() + numStatRegular() + numStatDevice() + numStatSymlink()
}

func numStatDirectory() uint64 {
	return uint64(len(statDirectory))
}

func numStatRegular() uint64 {
	return uint64(len(statRegular))
}

func numStatDevice() uint64 {
	return uint64(len(statDevice))
}

func numStatSymlink() uint64 {
	return uint64(len(statSymlink))
}

/*
func numStatUnsupported() uint64 {
	return uint64(len(statUnsupported))
}

func numStatInvalid() uint64 {
	return uint64(len(statInvalid))
}

func numStatIgnored() uint64 {
	return uint64(len(statIgnored))
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
	printStat(statDirectory, DIR_STR)
}

func printStatRegular() {
	printStat(statRegular, REG_STR)
}

func printStatDevice() {
	printStat(statDevice, DEVICE_STR)
}

func printStatSymlink() {
	printStat(statSymlink, SYMLINK_STR)
}
*/

func printStatUnsupported() {
	printStat(statUnsupported, UNSUPPORTED_STR)
}

func printStatInvalid() {
	printStat(statInvalid, INVALID_STR)
}

func printStatIgnored() {
	printStat(statIgnored, "ignored file")
}

func printStat(l []string, msg string) {
	if len(l) == 0 {
		return
	}
	printNumFormatString(uint64(len(l)), msg)

	for _, v := range l {
		f := getRealPath(v)
		t1, _ := getRawFileType(v)
		t2, _ := getFileType(v)
		assert(t2 != SYMLINK) // symlink chains resolved
		if t1 == SYMLINK {
			assert(optIgnoreSymlink || t2 == DIR || t2 == INVALID)
			fmt.Printf("%s (%s -> %s)\n",
				f, getFileTypeString(t1), getFileTypeString(t2))
		} else {
			assert(t2 != DIR)
			fmt.Printf("%s (%s)\n", f, getFileTypeString(t1))
		}
	}
}

// num written
func numWrittenTotal() uint64 {
	return numWrittenDirectory() + numWrittenRegular() + numWrittenDevice() + numWrittenSymlink()
}

func numWrittenDirectory() uint64 {
	return writtenDirectory
}

func numWrittenRegular() uint64 {
	return writtenRegular
}

func numWrittenDevice() uint64 {
	return writtenDevice
}

func numWrittenSymlink() uint64 {
	return writtenSymlink
}

// append written
func appendWrittenTotal(written uint64) {
}

func appendWrittenDirectory(written uint64) {
	writtenDirectory += written
}

func appendWrittenRegular(written uint64) {
	writtenRegular += written
}

func appendWrittenDevice(written uint64) {
	writtenDevice += written
}

func appendWrittenSymlink(written uint64) {
	writtenSymlink += written
}
