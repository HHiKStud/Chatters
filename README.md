"Chatters" is my attempt of creating a simple Chat Web-application with JWT tokens for authentication and web-sockets to establish connection. 
Made with Golang(and a little bit of JS on the front side).

All messages and users are stored in a local DB.

If you wanna try it out make sure you have PostgreDB installed on your machine.
Create a database chat_app there, and set a password 'pass' for default user(postgres).
You can configure everything in the config.go file.
Note that the password for user is mandatory, otherwise it won't work.
