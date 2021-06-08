Task2

In this we created three endpoints /login /signup /secretpage

Run the http_main.go in the terminal-1, then open other terminal-2 to and use curl "curl http://localhost:9000/endpoint"  this would prompt user and ask information in terminal-1 where we would give the required input. 


For Signup we take 3 input form user Username Roll number and password. We check if user has not alresdy signedup and also the inpu string are valid.Then we apply hashing and salting and then made an entry for the user in our database.

For login we take in the roll number and password and check if username and password matches . If so then we geerate a jwt token valid for 15 minutes. And login the user.

For secretpage we take in the jwt token from user.Now since the token is only present to the user logged in. So we could display the secret info to the user if the token is valid
