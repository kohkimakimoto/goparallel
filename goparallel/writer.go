package goparallel

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

var printMutex sync.Mutex

type prefixedWriter struct {
	Writer io.Writer
	Prefix string

	newline bool
}

func newPrefixWriter(writer io.Writer, prefix string) *prefixedWriter {
	return &prefixedWriter{
		Writer:  writer,
		Prefix:  prefix,
		newline: true,
	}
}

func (w *prefixedWriter) Write(data []byte) (int, error) {
	// in order to guarantee writing a single line by the one process. use mutex.
	printMutex.Lock()
	defer printMutex.Unlock()

	dataStr := string(data)
	dataStr = strings.Replace(dataStr, "\r\n", "\n", -1)

	var hasNewline = false
	if dataStr[len(dataStr)-1:] == "\n" {
		hasNewline = true
	}

	if w.Prefix != "" {
		if hasNewline {
			dataStr = strings.Replace(dataStr, "\n", "\n"+w.Prefix, strings.Count(dataStr, "\n")-1)
		} else {
			dataStr = strings.Replace(dataStr, "\n", "\n"+w.Prefix, -1) + "\n"
		}

		fmt.Fprint(w.Writer, w.Prefix+dataStr)
	} else {
		if hasNewline {
			// nothing to do.
		} else {
			dataStr = dataStr + "\n"
		}

		fmt.Fprint(w.Writer, dataStr)
	}

	return len(data), nil
}

//func (w *PrefixedWriter) Write(data []byte) (int, error) {
//	dataStr := string(data)
//	dataStr = strings.Replace(dataStr, "\r\n", "\n", -1)
//
//	if w.newline {
//		w.newline = false
//		fmt.Fprintf(w.Writer, "%s", w.Prefix)
//	}
//
//	if strings.Contains(dataStr, "\n") {
//		lineCount := strings.Count(dataStr, "\n")
//
//		if dataStr[len(dataStr)-1:] == "\n" {
//			w.newline = true
//		}
//
//		if w.newline {
//			dataStr = strings.Replace(dataStr, "\n", "\n"+w.Prefix, lineCount-1)
//		} else {
//			dataStr = strings.Replace(dataStr, "\n", "\n"+w.Prefix, -1)
//		}
//
//		fmt.Fprintf(w.Writer, "%s", dataStr)
//
//	} else {
//		fmt.Fprintf(w.Writer, "%s", dataStr)
//	}
//
//	return len(data), nil
//}
