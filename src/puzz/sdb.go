package puzz
import "sync"



// DB that wraps another DB and syncronizes calls.
type SyncDB struct {
        db    DB
        mutex *sync.Mutex
}

func (sdb *SyncDB) AddPicKV(k, v []byte) error {
        sdb.mutex.Lock()
        defer sdb.mutex.Unlock()
        return sdb.db.AddPicKV(k, v)
}

func (sdb *SyncDB) GetPicKV(k []byte) ([]byte, error) {
        sdb.mutex.Lock()
        defer sdb.mutex.Unlock()
        return sdb.db.GetPicKV(k)
}

func (sdb *SyncDB) GetNextId() (int, error) {
        sdb.mutex.Lock()
        defer sdb.mutex.Unlock()
        return sdb.db.GetNextId()
}

func (sdb *SyncDB) PutWord(k []byte) error {
        sdb.mutex.Lock()
        defer sdb.mutex.Unlock()
        return sdb.db.PutWord(k)
}


func (sdb *SyncDB) GetAllWords(prefix []byte) ([][]byte, error) {
        sdb.mutex.Lock()
        defer sdb.mutex.Unlock()
        return sdb.db.GetAllWords(prefix)
} 

