"Chatters" is my vision of a simple Chat Web-application with JWT tokens for authentication and web-sockets for establishing a connection between users. 
Made with Golang(and a little bit of JS on the front side).

All messages and users are stored in a local DB.

If you wanna try it out make sure you have PostgreSQL installed.
Create a database chat_app there, and set a password 'pass' for default user(postgres).
You can configure everything about DB in the config.go file if you want to.
Note that having a set password for user is mandatory, otherwise you won't be able to connect to DB.
