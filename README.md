# IITK Coin

I have split my package into multiple sub-packages (i.e. I have made a few sub-directories - `functions`, `global`, `handlers`).

Please look at the response format for each endpoint in the **`global`** package, where I have defined some variables and struct types, to be used in other functions.

In the three route handler functions I have used a `struct` (of types as in the `global` package) assigned to the variable named `payload`, for the body of the response to be sent. The three of them have similar structure. I have later converted these `struct`s into json objects.

In order to access /secretpage please provide an authorization token in the **header** in this format --> **"Authorization: Bearer \<access token\>"**.

Deliberately left `.env` out of `.gitignore` for the purposes of checking.
  
---
### For my reference:
- [X] look up popular directory structures
- [X] send a json object in the response
- [X] use MDN: HTTP status codes -> http.StatusBadRequest, http.StatusUnauthorized
- [ ] use refresh token
- [ ] A common approach for invalidating tokens when a user changes their password is to sign the token with a hash of their password. Thus if the password changes, any previous tokens automatically fail to verify. You can extend this to logout by including a last-logout-time in the user's record and using a combination of the last-logout-time and password hash to sign the token. This requires a DB lookup each time you need to verify the token signature, but presumably you're looking up the user anyway. – [Travis Terry](https://stackoverflow.com/questions/21978658/invalidating-json-web-tokens/23089839#comment45057142_23089839)
- [ ] (opt) check isAdmin from token and then authorize to /secretPage
- [ ] keep checking for unhandled errors
- [ ] a new table `auth` for storing just rollnos and passwords?
