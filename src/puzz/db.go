package puzz

import "fmt"
import "strconv"
import "sync"

var (
	// The DB we are going to use.
	db DB
	// Next id for new entries.
	next_id int
)

// DB provides access to two KV stores, one for pics and one for words.
type DB interface {
	// Add a key and a value to the Pic kv store.
	AddPicKV(k, v []byte) error

	// Get the value for a given key from the Pic kv store.
	GetPicKV(k []byte) ([]byte, error)

	// Loop over db and get the next id to use.
	GetNextId() (int, error)

	// insert a key into the word db.
	PutWord(k []byte) error

	// Get all the keys in the word db that start with prefix.
	GetAllWords(prefix []byte) ([][]byte, error)
}

func InitDb() error {
	var err error
	db, err = InitLDB()
	db = &SyncDB{
		db:    db,
		mutex: &sync.Mutex{},
	}
	if err != nil {
		return err
	}
	next_id, err = db.GetNextId()
	fmt.Println("DB intialized, next id:", next_id)
	return err
}

func SignatureKey(id int) []byte {
	return []byte(fmt.Sprintf("%09d", id))
}

func ParseSignatureKey(key []byte) (int, error) {
	if len(key) != 9 {
		return 0, fmt.Errorf("SignatureKey: wrong length %d", len(key))
	}
	id, err := strconv.Atoi(string(key))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SignatureValue(cvec []byte, name string) []byte {
	if got, want := len(cvec), 544; got != want {
		fmt.Printf("SignatureValue: wrong length: got %d; want %d", got, want)
	}
	return append(cvec, []byte(name)...)
}

func ParseSignatureValue(key []byte) ([]byte, string, error) {
	return key[:544], string(key[544:]), nil
}

func PartWordKey(word []byte, pos int) []byte {
	return append(word, []byte(fmt.Sprintf("%03d", pos))...)
}

func WordKey(word []byte, pos int, id int) []byte {
	return append(PartWordKey(word, pos), SignatureKey(id)...)
}

func ParseWordKey(key []byte) ([]byte, int, error) {
	id, err := ParseSignatureKey(key[K+3:])
	return key[:K], id, err
}

func AddPic(cvec []byte, name string) error {
	//fmt.Println("Adding", name, "with id", next_id)
	id := next_id
	next_id++
	key := SignatureKey(id)
	value := SignatureValue(cvec, name)
	// fmt.Printf("Adding signature key(%d) value(%d)\n", len(key), len(value))
	err := db.AddPicKV(key, value)
	if err != nil {
		return err
	}

	subvec := GetSubCVec(cvec)

	for pos, word := range subvec {
		key := WordKey(word, pos, id)
		// fmt.Printf("Adding word key(%d) value(%d)\n", len(key), 0)
		err = db.PutWord(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func HasPic(cvec []byte) ([]string, error) {
	candidates := make(map[int]int)
	for pos, word := range GetSubCVec(cvec) {
		partKey := PartWordKey(word, pos)
		keys, err := db.GetAllWords(partKey)
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			_, id, err := ParseWordKey(key)
			if err != nil {
				return nil, err
			}
			candidates[id]++
		}
	}
	result := []string{}
	for id, _ := range candidates {
		key := SignatureKey(id)
		value, err := db.GetPicKV(key)
		if len(value) == 0 {
			fmt.Println("pic not found", id)
		}
		canidate_cvec, name, err := ParseSignatureValue(value)
		if err != nil {
			return nil, err
		}

		delta := CompareCVecs(cvec, canidate_cvec)
		if delta < 0.6 {
			result = append(result, name)
		} else {
			fmt.Println("Found false candidate")
		}
	}
	return result, nil
}
