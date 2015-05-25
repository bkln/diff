package puzz

import "github.com/cznic/kv"
import "io"
import "bytes"

type KVDB struct {
	words, pics *kv.DB
}

func (kvdb *KVDB) AddPicKV(k, v []byte) error {
	err := kvdb.pics.Set(k, v)
	return err
}

func (kvdb *KVDB) GetPicKV(k []byte) ([]byte, error) {
	value, err := kvdb.pics.Get(nil, k)
	if len(value) == 0 {
		return []byte{}, nil
	}
	return value, err
}

func (kvdb *KVDB) GetNextId() (int, error) {
	next_id := 0
	iter, err := kvdb.pics.SeekFirst()
	if err == io.EOF {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	for {
		k, _, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		id, err := ParseSignatureKey(k)
		if err != nil {
			return 0, err
		}
		if id >= next_id {
			next_id = id + 1
		}
	}
	return next_id, nil
}

func (kvdb *KVDB) PutWord(k []byte) error {
	return kvdb.words.Set(k, []byte{})
}

func (kvdb *KVDB) GetAllWords(prefix []byte) ([][]byte, error) {
	result := make([][]byte, 0)
	iter, _, err := kvdb.words.Seek(prefix)
	if err != nil {
		return nil, err
	}
	for {
		k, _, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !bytes.HasPrefix(k, prefix) {
			break
		}
		kk := make([]byte, len(k))
		copy(kk, k)
		result = append(result, kk)
	}
	return result, nil

}

func InitKVDB() (*KVDB, error) {
	ops := kv.Options{}
	// database of [000000ID -> CVEC|picname]
	pics, err := kv.CreateTemp("./kdb", "pics", "", &ops)

	if err != nil {
		return nil, err
	}
	// database of [WORD|0000ID -> ]
	words, err := kv.CreateTemp("./kdb", "words", "", &ops)
	if err != nil {
		return nil, err
	}

	return &KVDB{
		pics:  pics,
		words: words,
	}, nil

}

func InitMemKVDB() (*KVDB, error) {
	ops := kv.Options{}
	// database of [000000ID -> CVEC|picname]
	pics, err := kv.CreateMem(&ops)

	if err != nil {
		return nil, err
	}
	// database of [WORD|0000ID -> ]
	words, err := kv.CreateMem(&ops)
	if err != nil {
		return nil, err
	}

	return &KVDB{
		pics:  pics,
		words: words,
	}, nil

}
