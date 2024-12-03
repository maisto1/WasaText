package database

const (
	usersTableCreationStatement = `
	CREATE TABLE "Users" (
		"user_id"			INTEGER NOT NULL UNIQUE,
		"username"			TEXT NOT NULL UNIQUE,
		"profile_photo"		BLOB,
		PRIMARY KEY("user_id" AUTOINCREMENT)
	);
	`
	conversationsTableCreationStatement = `
	CREATE TABLE "Conversations" (
		"conversation_id"	INTEGER NOT NULL UNIQUE,
		"name"				TEXT,
		"photo"				BLOB,
		"type"				TEXT,
		PRIMARY KEY("conversation_id" AUTOINCREMENT)
	);
	`

	messagesTableCreationStatement = `
	CREATE TABLE "Messages" (
		"message_id"		INTEGER NOT NULL UNIQUE,
		"conversation_id"	INTEGER NOT NULL,
		"user_id"			INTEGER NOT NULL,
		"content"			TEXT,
		"media"				BLOB,
		"type"				TEXT,
		"timestamp" 		INTEGER, 
		"status"			TEXT,
		"isForwarded" 		INTEGER,
		PRIMARY KEY("message_id" AUTOINCREMENT),
		FOREIGN KEY("conversation_id") REFERENCES "Conversations"("conversation_id") ON DELETE CASCADE,
		FOREIGN KEY("user_id") REFERENCES "Users"("user_id") ON DELETE CASCADE
	);
	`

	commentsTableCreationStatement = `
	CREATE TABLE "Comments" (
		"comment_id"		INTEGER NOT NULL UNIQUE,
		"message_id"		INTEGER NOT NULL,
		"user_id"			INTEGER NOT NULL,
		"content"			TEXT,
		"timestamp" 		INTEGER,
		PRIMARY KEY("comment_id" AUTOINCREMENT),
		FOREIGN KEY("message_id") REFERENCES "Messages"("message_id") ON DELETE CASCADE,
		FOREIGN KEY("user_id") REFERENCES "Users"("user_id") ON DELETE CASCADE
	);
	`

	partecipantsTableCreationStatement = `
	CREATE TABLE "Partecipants" (
		"user_id"          		INTEGER NOT NULL,
		"conversation_id"  		INTEGER NOT NULL,
		"conversation_name"		TEXT,
		"conversation_photo"	BLOB,
		PRIMARY KEY("user_id", "conversation_id"),
		FOREIGN KEY("user_id") REFERENCES "Users"("user_id") ON DELETE CASCADE,
		FOREIGN KEY("conversation_id") REFERENCES "Conversations"("conversation_id") ON DELETE CASCADE
	);
	`
)
