package tupleauth

import (
	"fmt"
	"testing"
)

func TestParseOkay(t *testing.T) {
	s := `
	family#member@1
	family#member@grandparents#member

	grandparents#member@2
	grandparents#member@3

	cats:pictures#read@family#member
	cats:pictures#write@1
	`

	records := Parse(s)
	if len(records) != 6 {
		t.Fatal("records len not as expected")
	}

	var (
		// users
		user1 = UserOrUserSet{UserID: 1}
		user2 = UserOrUserSet{UserID: 2}
		user3 = UserOrUserSet{UserID: 3}

		// objects
		family      = Object{ID: "family"}
		grandparent = Object{ID: "grandparents"}
		catPhotos   = Object{Namespace: "cats", ID: "pictures"}

		// user sets
		familyMember      = ObjectRelation{Object: family, Relation: "member"}
		grandparentMember = ObjectRelation{Object: grandparent, Relation: "member"}
	)

	expected := []Record{
		{Target: familyMember, Scope: user1},
		{Target: familyMember, Scope: UserOrUserSet{UserSet: grandparentMember}},

		{Target: grandparentMember, Scope: user2},
		{Target: grandparentMember, Scope: user3},

		{
			Target: ObjectRelation{Object: catPhotos, Relation: "read"},
			Scope:  UserOrUserSet{UserSet: familyMember},
		},
		{
			Target: ObjectRelation{Object: catPhotos, Relation: "write"},
			Scope:  user1,
		},
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
