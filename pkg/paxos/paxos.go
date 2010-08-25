package paxos

import (
	"os"
	"strings"
	"strconv"
)

type Msg struct {
	Seqn uint64
	From uint64
	To uint64
	Cmd string
	Body string
}

func (m Msg) SeqnX() uint64 {
	return m.Seqn
}

func (m Msg) FromX() uint64 {
	return m.From
}

func (m Msg) ToX() uint64 {
	return m.To
}

func (m Msg) CmdX() string {
	return m.Cmd
}

func (m Msg) BodyX() string {
	return m.Body
}

// TODO maybe we can make a better name for this. Not sure.
//
// SelfIndex is the position of the local node in the alphabetized list of all
// nodes in the cluster.
type Cluster interface {
	Putter
	Len() int
	Quorum() int
	SelfIndex() int
}

const (
	mFrom = iota
	mTo
	mCmd
	mBody
	mNumParts
)

var (
	InvalidArgumentsError = os.NewError("Invalid Arguments")
	Continue = os.NewError("continue")
)

func splitBody(body string, n int) ([]string, os.Error){
	bodyParts := strings.Split(body, ":", n)
	if len(bodyParts) != n {
		return nil, InvalidArgumentsError
	}
	return bodyParts, nil
}

func splitExactly(body string, n int) []string {
	parts, err := splitBody(body, n)
	if err != nil {
		panic(Continue)
	}
	return parts
}

func dtoui64(s string) uint64 {
	i, err := strconv.Btoui64(s, 10)
	if err != nil {
		panic(Continue)
	}
	return i
}

func swallowContinue() {
	p := recover()
	switch v := p.(type) {
	default: panic(p)
	case nil: return // didn't panic at all
	case os.Error:
		switch v {
		default: panic(v)
		case Continue: return
		}
	}
}

