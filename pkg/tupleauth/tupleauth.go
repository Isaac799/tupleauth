// Package tupleauth is inspired by https://authzed.com/zanzibar
package tupleauth

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	DelimiterObjectRelation  = "#"
	DelimiterRelationUser    = "@"
	DelimiterNamespaceObject = ":"
)

var (
	ErrWriteToDisk     = errors.New("failed writting to disk")
	ErrRestoreFromDisk = errors.New("failed restoring from disk")
	ErrRestoreEmpty    = errors.New("restore was empty")
)

type Object struct {
	Namespace string
	ID        string
}

func (obj *Object) IsZero() bool {
	return len(obj.Namespace)+len(obj.ID) == 0
}

func (obj *Object) String() string {
	if len(obj.ID) == 0 {
		return ""
	}
	if len(obj.Namespace) == 0 {
		return obj.ID
	}
	return strings.Join(
		[]string{
			obj.Namespace, DelimiterNamespaceObject,
			obj.ID,
		}, "")
}

type ObjectRelation struct {
	Object   Object
	Relation string
}

func (rel *ObjectRelation) IsZero() bool {
	return len(rel.Relation) == 0 || rel.Object.IsZero()
}

func (rel *ObjectRelation) String() string {
	if len(rel.Relation) == 0 {
		return ""
	}
	objStr := rel.Object.String()
	if len(objStr) == 0 {
		return ""
	}
	return strings.Join(
		[]string{
			objStr, DelimiterObjectRelation,
			rel.Relation,
		}, "")
}

type UserOrUserSet struct {
	UserID int

	// below unused if UserID > 0
	UserSet ObjectRelation
}

func (usr *UserOrUserSet) IsZero() bool {
	return usr.UserID == 0 && usr.UserSet.IsZero()
}

func (usr *UserOrUserSet) String() string {
	if usr.UserID > 0 {
		return strconv.Itoa(usr.UserID)
	}
	return usr.UserSet.String()
}

type Record struct {
	Iat time.Time

	Target ObjectRelation
	Scope  UserOrUserSet
}

func (rec *Record) IsZero() bool {
	return rec.Iat.IsZero() || rec.Target.IsZero() || rec.Scope.IsZero()
}

func (rec *Record) String() string {
	return strings.Join(
		[]string{
			rec.Target.String(), DelimiterRelationUser,
			rec.Scope.String(),
		}, "")
}
