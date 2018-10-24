package knife

import (
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

const (
	BUF_LEN = 32768
	SUTFF   = "part"
)

type Chunk struct {
	off int64
	buf []byte
}

func CutBySize(filename string, size int64) {
	fs, e := os.Stat(filename)
	if nil != e {
		panic(e)
	}
	fsize := fs.Size()
	if fsize < size {
		panic(fmt.Errorf("invalid size"))
	}
	n := int64(math.Ceil(float64(fsize) / float64(size)))

	e = SavePartInfo(fmt.Sprintf("%s.partinfo", filename),
		PartInfo{Filename: filename, Parts: n})
	if e != nil {
		panic(e)
	}

	if n > 1 {
		wg := &sync.WaitGroup{}
		var i int64
		for i = 0; i < n-1; i++ {
			go func(i int64) {
				wg.Add(1)
				defer wg.Done()
				cut(fmt.Sprintf("%s.%s%v", filename, SUTFF, i), filename, size*i, int(size))
			}(i)
		}
		go func(i int64) {
			wg.Add(1)
			defer wg.Done()
			cut(fmt.Sprintf("%s.%s%v", filename, SUTFF, i), filename, size*i, int(fsize-size*i))
		}(n - 1)
		wg.Wait()
	} else {
		cut(fmt.Sprintf("%s.%s%v", filename, SUTFF, 0), filename, 0, int(fsize))
	}
}

func cut(dst, src string, off int64, size int) {
	e := os.Remove(dst)
	if nil != e && !os.IsNotExist(e) {
		panic(e)
	}

	fdst, e := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.ModePerm)
	if nil != e {
		panic(e)
	}
	defer fdst.Close()

	fsrc, e := os.Open(src)
	if nil != e {
		panic(e)
	}
	defer fsrc.Close()
	if _, e = fsrc.Seek(off, 0); e != nil {
		panic(e)
	}

	chunkCh := make(chan Chunk)
	go func() {
		reader(fsrc, size, BUF_LEN, chunkCh, 0)
		close(chunkCh)
	}()

	for {
		select {
		case chunk, ok := <-chunkCh:
			if !ok {
				return
			}
			fdst.WriteAt(chunk.buf, chunk.off)
		}
	}
}

func reader(f *os.File, size int, buflen int, chunkCh chan Chunk, off int64) {
	if size < buflen {
		buflen = size
	}
	left := size
	for {
		chunk := Chunk{
			off: off,
			buf: make([]byte, buflen),
		}
		n, e := f.Read(chunk.buf)
		if e != nil {
			if e == io.EOF {
				return
			}
			panic(e)
		}
		chunkCh <- chunk
		left -= n
		off += int64(n)
		if left <= 0 {
			return
		} else if left < buflen {
			buflen = left
		}
	}
}

func Pack(filename, outputname string) {
	info, e := LoadPartInfo(filename)
	if e != nil {
		panic(e)
	}

	if len(outputname) == 0 {
		outputname = info.Filename
	}
	fdst, e := os.OpenFile(outputname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if e != nil {
		panic(e)
	}
	defer fdst.Close()

	chunkCh := make(chan Chunk)

	count := info.Parts
	var off int64 = 0
	var i int64
	for i = 0; i < info.Parts; i++ {
		fname := fmt.Sprintf("%s.%s%d", info.Filename, SUTFF, i)
		fi, e := os.Stat(fname)
		if e != nil {
			panic(e)
		}
		f, e := os.Open(fname)
		if e != nil {
			panic(e)
		}
		defer f.Close()
		go func(off int64) {
			reader(f, int(fi.Size()), BUF_LEN, chunkCh, off)
			count -= 1
			if count <= 0 {
				close(chunkCh)
			}
		}(off)
		off += fi.Size()
	}

	for {
		select {
		case chunk, ok := <-chunkCh:
			if !ok {
				return
			}
			fdst.WriteAt(chunk.buf, chunk.off)
		}
	}
}
