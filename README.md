# IITK Coin

I have split my package into multiple sub-packages (i.e. I have made a few sub-directories - `functions`, `global`, `handlers`).

For the response format for each endpoint please look at the **`global`** package, where I have defined some variables and struct types, to be used in other functions.

The requests format for each endpoint is specified in the beginning of each of the handler functions.

Although the endpoints have slightly different format for their response object, I have handled them all using a `type-switch` in a common `server.Respond` function which responds to requests for all the endpoints.

In order to access /secretpage please provide an authorization token in the **header** in this format --> **"Authorization: Bearer \<access token\>"**.

Deliberately left `.env` out of `.gitignore` for the purposes of checking.
  
---
### For my reference:
- [x] look up popular directory structures
- [x] send a json object in the response
- [x] use MDN: HTTP status codes -> http.StatusBadRequest, http.StatusUnauthorized
- [ ] use refresh token
- [ ] A common approach for invalidating tokens when a user changes their password is to sign the token with a hash of their password. Thus if the password changes, any previous tokens automatically fail to verify. You can extend this to logout by including a last-logout-time in the user's record and using a combination of the last-logout-time and password hash to sign the token. This requires a DB lookup each time you need to verify the token signature, but presumably you're looking up the user anyway. – [Travis Terry](https://stackoverflow.com/questions/21978658/invalidating-json-web-tokens/23089839#comment45057142_23089839)
- [ ] (opt) check isAdmin from token and then authorize to /secretPage
- [ ] keep checking for unhandled errors
- [ ] a new table `auth` for storing just rollnos and passwords?
