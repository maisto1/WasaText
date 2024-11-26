package database

const (
	usersTableCreationStatement = `
	CREATE TABLE Users (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Username VARCHAR(255) NOT NULL,
		ProfilePhoto TEXT DEFAULT NULL
	);
	`
	conversationsTableCreationStatement = `
	CREATE TABLE Conversations (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Name VARCHAR(255) NOT NULL,
		ConversationType TEXT NOT NULL,
		ConversationPhoto TEXT DEFAULT NULL
	);
	`

	userConversationsTableCreationStatement = `
	CREATE TABLE UserConversations (
		UserID INT NOT NULL,
		ConversationID INT NOT NULL,
		LastMessageContent TEXT DEFAULT NULL,
		LastMessageTimestamp DATETIME DEFAULT NULL,
		PRIMARY KEY (UserID, ConversationID),
		FOREIGN KEY (UserID) REFERENCES Users(ID) ON DELETE CASCADE,
		FOREIGN KEY (ConversationID) REFERENCES Conversations(ID) ON DELETE CASCADE
	);
	`

	conversationsPartecipantTableCreationStatement = `
	CREATE TABLE ConversationParticipants (
		ConversationID INT NOT NULL,
		UserID INT NOT NULL,
		PRIMARY KEY (ConversationID, UserID),
		FOREIGN KEY (ConversationID) REFERENCES Conversations(ID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES Users(ID) ON DELETE CASCADE
	);
	`

	messagesTableCreationStatement = `
	CREATE TABLE Messages (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		ConversationID INT NOT NULL,
		SenderID INT NOT NULL,
		Timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		MediaType VARCHAR(50) NOT NULL,
		MessageContent TEXT DEFAULT NULL,
		Photo TEXT DEFAULT NULL,
		Status VARCHAR(50) NOT NULL,
		IsForwarded BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (ConversationID) REFERENCES Conversations(ID) ON DELETE CASCADE,
		FOREIGN KEY (SenderID) REFERENCES Users(ID) ON DELETE CASCADE
	);
	`

	commentsTableCreationStatement = `
	CREATE TABLE Comments (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		MessageID INT NOT NULL,
		UserID INT NOT NULL,
		CommentContent TEXT NOT NULL,
		Timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (MessageID) REFERENCES Messages(ID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES Users(ID) ON DELETE CASCADE
	);
	`
)
