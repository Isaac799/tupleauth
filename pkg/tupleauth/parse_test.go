package tupleauth

import (
	"fmt"
	"testing"
)

func TestParseOkay(t *testing.T) {
	s := `
	family#member@john
	family#member@grandparents#member

	grandparents#member@linda
	grandparents#member@larry

	cats:pictures#read@family#member
	cats:pictures#write@john
	`
	records := Parse(s)
	if len(records) != 6 {
		t.Fatal("records len not as expected")
	}

	var (
		// users
		john  = UserSet{UserID: "john"}
		larry = UserSet{UserID: "larry"}
		linda = UserSet{UserID: "linda"}

		// objects
		family       = Object{ObjectID: "family"}
		grandparents = Object{ObjectID: "grandparents"}
		catPhotos    = Object{Namespace: "cats", ObjectID: "pictures"}

		// user set
		familyMember   = UserSet{Object: family, Relation: "member"}
		grandparentRef = UserSet{Object: grandparents, Relation: "member"}
	)

	expected := []Record{
		{Obj: family, Rel: "member", Usr: john},
		{Obj: family, Rel: "member", Usr: grandparentRef},

		{Obj: grandparents, Rel: "member", Usr: linda},
		{Obj: grandparents, Rel: "member", Usr: larry},

		{Obj: catPhotos, Rel: "read", Usr: familyMember},
		{Obj: catPhotos, Rel: "write", Usr: john},
	}

	for i, expect := range expected {
		actual := records[i]
		expectStr := expect.String()
		actualStr := actual.String()
		if actualStr != expectStr {
			s := fmt.Sprintf("EXPECT != ACTUAL \nEXPECT: '%s'\nACTUAL: '%s'", expectStr, actualStr)
			t.Fatal(s)
		}
	}
}
