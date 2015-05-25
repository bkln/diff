package puzz

import "net/http"
import "fmt"

type DupOrInsertArg struct {
	Name string
	Cvec []byte
}

type DupOrInsertReply struct {
	Dups []string
}

type DupOrInsertService struct{}

func (d *DupOrInsertService) DupOrInsert(r *http.Request, args *DupOrInsertArg, reply *DupOrInsertReply) error {
	if len(args.Cvec) != 544 {
		return fmt.Errorf("Wrong length, not 544")
	}

	duplicates, err := HasPic(args.Cvec)
	if err != nil {
		return err
	}
	if len(duplicates) == 0 {
		err := AddPic(args.Cvec, args.Name)
		if err != nil {
			return err
		}
	}
	reply.Dups = duplicates
	return nil
}
