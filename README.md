# IITK Coin
This is a `go` program. Once executed the program connects to a SQLite database named `iitkusers` and creates a table named `users` inside the database. It contains a function that takes details from a user as its arguement and adds them to the database only if they are a new user. If however, the details of an existing user are tried to be added, the function just displays that the user details are already present and does nothing more.