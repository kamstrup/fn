Tips and Tricks
====

Limit Use of Inline Function Bodies
----
It is very tempting to use inline functions a lot with doing functional programming,
but it makes the code harder to read, harder to test, and harder to reuse.
It is generally good practice in functional programming to have short and concise
method bodies. This applies when using Fn as well.

Instead of this
```.go
userIDs := seq.SliceOfArgs("user1", "user2", ...)
users := seq.Mapping(userIDs, func(userID string) opt.Opt[*User] {
    db, err := getDBConnection()
    if err != nil {
        return opt.ErrorOf[*User](err)
    }
    user, err := db.GetUser(userID)
    if err != nil {
        return opt.ErrorOf[*User](err)
    }
    return opt.Of(user)
}).Where(opt.Ok[*User])
```

Do this: Put the body of the mapping function into a separate function:
```.go
func findUser(userID string) opt.Opt[*User] {
    db, err := getDBConnection()
    ...
}
```
and your user lookups become
```.go
userIDs := seq.SliceOfArgs("user1", "user2", ...)
users := seq.Mapping(userIDs, findUser).
    Where(opt.Ok[*User])
```

Mapping Methods
----
You can use references to methods where a function is required.
Let's explore with an example.

In our previous example the function `findUser()` retrieved the DB connection
on every call. It would be cleaner object-oriented design if we moved that function
onto the DB connection as a method, and the DB API might not want to expose the opt
package in the public API. That would look like:
```.go
package mydb

func (db *DBConnection) FindUser(userID string) (*User, error) {
   ...
}
```
To fix our user lookup seq from before we can write it like this now:
```.go
db := getDBConnection()
userIDs := seq.SliceOfArgs("user1", "user2", ...)
users := seq.Mapping(userIDs, opt.Mapper(db.FindUser)).
    Where(opt.Ok[*User])
```

Method Expressions
----
Consider a package `mydb` with a `User` struct with an `ID()` method
```.go
package mydb

func (u *User) ID() string {
   ...
}
```
If we want to convert a slice of users into a slice of user IDs we could do
```.go
users := seq.SliceOf(db.GetAllUsers())
userIDs := seq.Mapping(users, func(u *mydb.User) string {
    return u.ID()
}).ToSlice()
```
But there is a shorter and better performing way to write this, by using "method expressions".
In Go you can grab a reference to a method like `*User.ID()` and turn it into a
function taking the method receiver as first argument `func (*mydb.User) string`:
```.go
users := seq.SliceOf(db.GetAllUsers())
userIDs := seq.Mapping(users, (*mydb.User).ID).ToSlice()
```
