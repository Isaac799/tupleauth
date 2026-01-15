package tupleauth

import (
	"strconv"
	"strings"
	"time"
)

func ParseUserSet(s string) *UserOrUserSet {
	before, after, found := strings.Cut(s, DelimiterObjectRelation)
	if !found {
		n, _ := strconv.Atoi(before)
		return &UserOrUserSet{
			UserID: n,
		}
	}

	return &UserOrUserSet{
		UserSet: ObjectRelation{
			Object:   *ParseObject(before),
			Relation: after,
		},
	}
}

func ParseObject(s string) *Object {
	before, after, found := strings.Cut(s, DelimiterNamespaceObject)
	if !found {
		return &Object{
			ID: before,
		}
	}
	return &Object{
		Namespace: before,
		ID:        after,
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
			Target: ObjectRelation{
				Object:   *ParseObject(object),
				Relation: relation,
			},
			Scope: *ParseUserSet(after),
		}
		if r.IsZero() {
			continue
		}
		records = append(records, r)
	}
	return records
}
