// Package tupleauth is inspired by https://authzed.com/zanzibar
package tupleauth

import (
	"errors"
	"strings"
	"time"
)

const (
	DelimiterObjectRelation  = "#"
	DelimiterRelationUser    = "@"
	DelimiterNamespaceObject = ":"
	DelimiterUserset         = "#"
)

var (
	ErrWriteToDisk     = errors.New("failed writting to disk")
	ErrRestoreFromDisk = errors.New("failed restoring from disk")
	ErrRestoreEmpty    = errors.New("restore was empty")
)

type Object struct {
	Namespace string
	ObjectID  string
}

func (r *Object) IsZero() bool {
	return len(r.Namespace)+len(r.ObjectID) == 0
}

func (r *Object) String() string {
	if len(r.ObjectID) == 0 {
		return ""
	}
	if len(r.Namespace) == 0 {
		return r.ObjectID
	}
	return strings.Join(
		[]string{
			r.Namespace, DelimiterNamespaceObject,
			r.ObjectID,
		}, "")
}

type UserSet struct {
	UserID string

	// below are ignored if UserID has len
	Object   Object
	Relation string
}

func (us *UserSet) IsZero() bool {
	return len(us.UserID) == 0 && us.Object.IsZero()
}

func (us *UserSet) String() string {
	if len(us.UserID) > 0 {
		return us.UserID
	}
	if len(us.Relation) == 0 {
		return ""
	}
	objStr := us.Object.String()
	if len(objStr) == 0 {
		return ""
	}
	return strings.Join(
		[]string{
			objStr, DelimiterUserset,
			us.Relation,
		}, "")
}

type Record struct {
	Iat time.Time
	Obj Object
	Rel string
	Usr UserSet
}

func (r *Record) IsZero() bool {
	return r.Iat.IsZero() || r.Obj.IsZero() || len(r.Rel) == 0 || r.Usr.IsZero()
}

func (r *Record) String() string {
	return strings.Join(
		[]string{
			r.Obj.String(), DelimiterObjectRelation,
			r.Rel, DelimiterRelationUser,
			r.Usr.String(),
		}, "")
}
