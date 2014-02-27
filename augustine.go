package augustine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
)

const (
	boltDBName = "bolt.db"

	// Bucket names
	bDoc = "docs"
)

type DB struct {
	b          *bolt.DB
	tokenizer  Tokenizer
	normalizer Normalizer
	ftsIndex   *FTSIndex

	dir     string
	nextSeq uint64
}

func (d *DB) createBoltBuckets() {
	for _, name := range []string{bDoc} {
		// Delete the bucket first to clear it
		_ = d.b.DeleteBucket(name)
		if err := d.b.CreateBucket(name); err != nil {
			panic(err) // Shouldn't happen -- only a bad bucket name causes this
		}
	}
}

func Open(dir string) (*DB, error) {
	stat, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	case !stat.IsDir():
		return nil, fmt.Errorf("%s is not a directory", dir)
	}
	var b bolt.DB
	if err := b.Open(filepath.Join(dir, boltDBName), 0666); err != nil {
		return nil, err
	}
	db := &DB{
		dir:        dir,
		b:          &b,
		tokenizer:  SimpleTokenizer,
		normalizer: NewSimpleNormalizer(),
		ftsIndex:   NewFTSIndex(),
	}
	db.createBoltBuckets()
	return db, nil
}

func (d *DB) Close() {
	d.b.Close()
}

// nextID increments and returns the sequence number
func (d *DB) nextID() key { return key(atomic.AddUint64(&d.nextSeq, 1)) }

func (d *DB) Put(doc *Doc) (id uint64, err error) {
	key := d.nextID()
	val, err := doc.MarshalBinary()
	if err != nil {
		return 0, err
	}
	if err := d.b.Put(bDoc, key.Bytes(), val); err != nil {
		return 0, err
	}
	d.index(key, doc)
	return uint64(key), nil
}

func (d *DB) Get(id uint64) (doc *Doc, err error) {
	key := key(id)
	val, err := d.b.Get(bDoc, key.Bytes())
	if err != nil {
		return nil, err
	}
	doc = new(Doc)
	if err := doc.UnmarshalBinary(val); err != nil {
		return nil, err
	}
	return doc, nil
}

// index adds doc's Text field to the search index.
func (d *DB) index(id key, doc *Doc) {
	tokenReader := d.tokenizer(bytes.NewBuffer(doc.Text))
	var pos int64
	for {
		tok, err := tokenReader.ReadToken()
		fmt.Printf("\033[01;34m>>>> tok: %v\x1B[m\n", tok)
		fmt.Printf("\033[01;34m>>>> err: %v\x1B[m\n", err)
		if err != nil {
			// Right now this can only arise from EOF.
			return
		}
		posting := Posting{
			docID: id,
			pos:   pos,
		}
		for _, term := range d.normalizer.Normalize(tok) {
			d.ftsIndex.Add(string(term), posting)
		}
		pos++
	}
}

func (d *DB) Print() {
	fmt.Println("Index:")
	spew.Printf("%#v\n", d.ftsIndex)
}
