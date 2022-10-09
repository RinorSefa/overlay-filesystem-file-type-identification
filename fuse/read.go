// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import (
	"io"
	"syscall"

	"github.com/gabriel-vasile/mimetype"
)

// ReadResultData is the read return for returning bytes directly.
type readResultData struct {
	// Raw bytes for the read.
	Data []byte
}

func (r *readResultData) Size() int {
	return len(r.Data)
}

func (r *readResultData) Done() {
}

func (r *readResultData) Bytes(buf []byte) ([]byte, Status) {
	return r.Data, OK
}

func ReadResultData(b []byte) ReadResult {
	return &readResultData{b}
}

func ReadResultFd(fd uintptr, off int64, sz int) ReadResult {
	return &readResultFd{fd, off, sz}
}

// ReadResultFd is the read return for zero-copy file data.
type readResultFd struct {
	// Splice from the following file.
	Fd uintptr

	// Offset within Fd, or -1 to use current offset.
	Off int64

	// Size of data to be loaded. Actual data available may be
	// less at the EOF.
	Sz int
}

// Author @RinorSefa
// mimetype/supported_mimes.md
// according to paper CryptoLock Scaife 2016 frequency of file extension accessed by 492 ransomware samples
// pdf / "application/pdf" / offset 0 (according to https://github.com/file/file/blob/master/magic/Magdir/pdf)
// odt / "application/vnd.oasis.opendocument.text" / (offset 50, 0 should work, as 0 + 4026)
// docx / "application/vnd.openxmlformats-officedocument.wordprocessingml.document" / offset 0
// pptx / "application/vnd.openxmlformats-officedocument.presentationml.presentation" / offset 0
// txt / "text/plain" / plain bytes/ nothing special
// mov / "video/quicktime" // no info.
// zip / "application/zip" // offset 0
// jpg / "image/jpeg" "0
// xls / "application/vnd.ms-excel" // 0
// csv / "text/csv" // file/file no description, special case checks into the buffer
// doc / "application/msword" (0 or 8, but both fall inside 0)
// png / "image/png" (0-12)
// file type; what stakeholders want to protect. GG WP
var typesList = [12]string{
	"application/pdf",                         //pdf
	"application/vnd.oasis.opendocument.text", // odt
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document", // docx
	"application/msword", // doc
	"application/vnd.openxmlformats-officedocument.presentationml.presentation", //pptx
	"application/vnd.ms-powerpoint",                                             //ppt
	"application/zip",                                                           //zip
	"image/jpeg",                                                                // jpeg
	"application/vnd.ms-excel",                                                  //xls
	"image/png",                                                                 // png
	"text/csv",                                                                  // csv
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", // .xlsx
}

// Modified by @RinorSefa
// Reads raw bytes from file descriptor if necessary, using the passed
// buffer as storage.
func (r *readResultFd) Bytes(buf []byte) ([]byte, Status) {
	sz := r.Sz
	if len(buf) < sz {
		sz = len(buf)
	}

	n, err := syscall.Pread(int(r.Fd), buf[:sz], r.Off)
	if err == io.EOF {
		err = nil
	}

	//our code
	if r.Off == 0 {
		mtype := mimetype.Detect(buf[:n])
		for _, v := range typesList {
			if mtype.Is(v) {
				buf = make([]byte, n)
				break
			}
		}
	}

	if n < 0 {
		n = 0
	}

	return buf[:n], ToStatus(err)
}

func (r *readResultFd) Size() int {
	return r.Sz
}

func (r *readResultFd) Done() {
}
