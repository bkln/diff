package puzz

import "fmt"
import "github.com/syndtr/goleveldb/leveldb"
import "github.com/syndtr/goleveldb/leveldb/opt"
import "github.com/syndtr/goleveldb/leveldb/util"

type LDB struct {
	words, pics *leveldb.DB
}

func (ldb *LDB) AddPicKV(k, v []byte) error {
	err := ldb.pics.Put(k, v, nil)
	return err
}

func (ldb *LDB) GetPicKV(k []byte) ([]byte, error) {
	value, err := ldb.pics.Get(k, nil)
	if err == leveldb.ErrNotFound {
		return []byte{}, nil
	}
	return value, err
}

func (ldb *LDB) GetNextId() (int, error) {
	next_id := 0
	iter := ldb.pics.NewIterator(nil, nil)
	for iter.Next() {
		id, err := ParseSignatureKey(iter.Key())
		if err != nil {
			return 0, err
		}
		if id >= next_id {
			next_id = id + 1
		}
	}
	iter.Release()
	return next_id, iter.Error()
}

func (ldb *LDB) PutWord(k []byte) error {
	return ldb.words.Put(k, []byte{}, nil)
}

func (ldb *LDB) GetAllWords(prefix []byte) ([][]byte, error) {
	result := make([][]byte, 0)
	iter := ldb.words.NewIterator(util.BytesPrefix(prefix), nil)
	for iter.Next() {
		k := make([]byte, len(iter.Key()))
		copy(k, iter.Key())
		result = append(result, k)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	}
	return result, nil

}

func InitLDB() (*LDB, error) {
	o := &opt.Options{
		BlockCacheCapacity: 1 * 1024*1024*1024,
	}
	fmt.Println("122")
	// database of [000000ID -> CVEC|picname]
	pics, err := leveldb.OpenFile("./db/pics", o)
	if err != nil {
		return nil, err
	}
	fmt.Println("111")
	// database of [WORD|0000ID -> ]
	words, err := leveldb.OpenFile("./db/words", o)
	if err != nil {
		return nil, err
	}

	return &LDB{
		pics:  pics,
		words: words,
	}, nil
}
