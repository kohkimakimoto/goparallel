package ltsv

// This code inspired by https://github.com/ymotongpoo/goltsv

//Copyright (c) 2012, Yoshifumi YAMAGUCHI (ymotongpoo)
//All rights reserved.
//
//Redistribution and use in source and binary forms, with or without modification,
//are permitted provided that the following conditions are met:
//
//* Redistributions of source code must retain the above copyright notice, this list
//of conditions and the following disclaimer.
//
//* Redistributions in binary form must reproduce the above copyright notice, this
//list of conditions and the following disclaimer in the documentation and/or other
//materials provided with the distribution.
//
//* Neither the name of the betterForm Project
//nor the names of its contributors may be used to endorse or promote products derived
//from this software without specific prior written permission.
//
//THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY
//EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
//OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT
//SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
//INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
//PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
//INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
//LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// These are the errors that can be returned in ParseError.Error
var (
	ErrFieldFormat = errors.New("wrong LTSV field format")
	ErrLabelName   = errors.New("unexpected label name")
)

// A Reader reads from a Labeled TSV (LTSV) file.
//
// As returned by NewReader, a Reader expects input conforming LTSV (http://ltsv.org/)
type LTSVReader struct {
	reader *bufio.Reader
}

// NewReader returns a new LTSVReader that reads from r.
func NewReader(r io.Reader) *LTSVReader {
	return &LTSVReader{bufio.NewReader(r)}
}

// error creates a new Error based on err.
func (r *LTSVReader) error(err error) error {
	return err
}

// Read reads one record from r. The record is a map of string with
// each key and value representing one field.
func (r *LTSVReader) Read() (record map[string]string, err error) {
	var line []byte
	record = make(map[string]string)

	for {
		line, _, err = r.reader.ReadLine()
		if err != nil {
			return nil, err
		}

		sline := strings.TrimSpace(string(line))
		if sline == "" {
			// Skip empty line
			continue
		}
		tokens := strings.Split(sline, "\t")
		if len(tokens) == 0 {
			return nil, r.error(ErrFieldFormat)
		}
		for _, field := range tokens {
			if field == "" {
				continue
			}
			data := strings.SplitN(field, ":", 2)
			if len(data) != 2 {
				return record, r.error(ErrLabelName)
			}
			record[data[0]] = data[1]
		}
		return record, nil
	}
	return
}

// ReadAll reads all the remainig records from r.
// Each records is a slice of map of fields.
func (r *LTSVReader) ReadAll() (records []map[string]string, err error) {
	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}
