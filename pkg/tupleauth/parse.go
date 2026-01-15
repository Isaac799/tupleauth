package tupleauth

import (
	"strings"
	"time"
)

func ParseUserSet(s string) *UserSet {
	before, after, found := strings.Cut(s, DelimiterUserset)
	if !found {
		return &UserSet{
			UserID: before,
		}
	}

	return &UserSet{
		Relation: after,
		Object:   *ParseObject(before),
	}
}

func ParseObject(s string) *Object {
	before, after, found := strings.Cut(s, DelimiterNamespaceObject)
	if !found {
		return &Object{
			ObjectID: before,
		}
	}
	return &Object{
		Namespace: before,
		ObjectID:  after,
	}
}

func Parse(s string) []Record {
	lines := strings.Lines(s)
	var records []Record
	for line := range lines {
		line = strings.TrimSpace(line)
		before, after, found := strings.Cut(line, DelimiterRelationUser)
		if !found {
			continue
		}
		object, relation, found := strings.Cut(before, DelimiterObjectRelation)
		if !found {
			continue
		}
		r := Record{
			Iat: time.Now(),
			Obj: *ParseObject(object),
			Rel: relation,
			Usr: *ParseUserSet(after),
		}
		if r.IsZero() {
			continue
		}
		records = append(records, r)
	}
	return records
}
