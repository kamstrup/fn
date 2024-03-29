package examples

import (
	"fmt"
	"testing"

	"github.com/kamstrup/fn/seq"
)

type User struct {
	id   string
	name string
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) HasName() bool {
	return len(u.name) != 0
}

func (u *User) String() string {
	return fmt.Sprintf("%s(%s)", u.name, u.id)
}

func TestExampleUsersById(t *testing.T) {
	// Assume we have collected a slice of Users from somewhere, a DB query maybe
	users := []*User{
		{
			id:   "jh234dg",
			name: "mikkel",
		},
		{
			id:   "8o4355",
			name: "",
		},
		{
			id:   "wegr98u2",
			name: "",
		},
		{
			id:   "h38ufg",
			name: "bob",
		},
		{
			id:   "x985ng",
			name: "mikkel",
		},
	}

	// Let's check if they all have a valid "name" field
	everyOneHasName := seq.All(seq.SliceOf(users), (*User).HasName)
	fmt.Println("Does everyone have a name?", everyOneHasName)

	// Now let's print the IDs of the users without names, sorted reverse alphabetically
	usersWithEmptyNames := seq.SliceOf(users).
		Where(seq.Not((*User).HasName))
	idsWithEmptyNames := seq.MappingOf(usersWithEmptyNames, (*User).ID).
		ToSlice().
		Sort(seq.OrderDesc[string])
	fmt.Println("These user IDs do not have a name:", idsWithEmptyNames)

	// Let's create a map[userID]*User:
	// First we create a Seq of Tuples(userId, User)
	usersWithIDs := seq.MappingOf(seq.SliceOf(users), seq.TupleWithKey((*User).ID))
	// Now flush that Seq of tuples into the MakeMap collector
	usersByIDs := seq.Reduce(seq.MakeMap[string, *User], nil, usersWithIDs).Or(nil)

	// usersById is now a map[string]*User. Let's look up some users
	fmt.Println("User with ID(xyz123):", usersByIDs["xyz123"]) // no one, nil
	fmt.Println("User with ID(h38ufg):", usersByIDs["h38ufg"]) // bob

	// Let's pretend that some new users come in as a map[userID]*User
	newUsers := map[string]*User{
		"987hgj": {
			id:   "987hgj",
			name: "beatrice",
		},
		"mnb456": {
			id:   "mnb456",
			name: "lena",
		},
	}

	// Let's get a combined list of Users sorted by name
	allUsers := seq.ConcatOf(seq.MapOf(usersByIDs), seq.MapOf(newUsers)) // MapOf handles maps as Seqs of Tuples
	allUsersSorted := seq.MappingOf(allUsers, seq.TupleValue[string, *User]).
		ToSlice().
		Sort(func(u1, u2 *User) bool { return u1.name < u2.name })

	fmt.Println("Combined list of users, by name:", allUsersSorted)

	// We have decided that users must have a non-empty name,
	// and the name must start with "a" or "b". Controversial, but here we are.
	fmt.Println("Users starting with 'a' or 'b':")
	allUsersSorted.
		Where(func(u *User) bool { return len(u.name) != 0 }). // only valid names
		While(func(u *User) bool { return u.name[0] <= 'b' }). // stop the Seq after 'b'
		ForEachIndex(func(i int, u *User) {
			fmt.Println("User", i, "is", u) // just beatrice and bob
		})
}
