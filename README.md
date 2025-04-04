"Chatters" is my vision of a simple Chat Web-application with JWT tokens for authentication and web-sockets for establishing a connection between users. 
Made with Golang(and a little bit of JS on the front side).

All messages and users are stored in a local DB.

If you wanna try it out make sure you have PostgreSQL installed.
Create a database chat_app there, and set a password 'pass' for default user(postgres).
You can configure everything about DB in the config.go file if you want to.
Note that having a set password for user is mandatory, otherwise you won't be able to connect to DB.

This applicatioin uses the following design patterns(even though they are not "clear"):
1. Singleton. For declaring structures and types in go.
2. Publisher-Subscriber. Websocket mechanism. Server(publisher) sends messages to all the clients(subscribers)
3. Facade. HTTP handlers in the main.go that hide all the code inside them.
4. Observer. For connecting/disconnect users from the chat.
5. Front Controller. For http-routing in main.go. All the handlers go through http.handleFunc.