/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"

	"github.com/maisto1/WasaText/service/models"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Login or Register
	Login(username string) (int64, error)

	// Get preview conversations
	GetPreviewConversations(user_id int64) ([]models.Preview, error)

	// Get Latest Conversation Message
	GetLatestMessage(conversation_id int64) (models.Message, error)

	// Create a new conversation
	CreateConversation(user_id int64, group_name string, typeConv string, partecipant string) (int, error)

	// Checks if 2 user already have a private conversation
	CheckPrivateConversation(user_id1, user_id2 int64) (bool, error)

	// Return every messages from a specific conversation
	GetMessages(user_id int64, conversation_id int64) ([]models.Message, error)

	// Send a message in a conversation
	CreateMessage(user_id int64, conversation_id int64, target_id int64, typeMessage string, content string, media []byte, forwarded bool) (models.Message, error)

	// Delete a message
	DeleteMessage(user_id int64, conversation_id int64, message_id int64) error

	// Forward a message to another conversation
	ForwardMessage(user_id int64, conversation_id int64, target_id int64, message_id int64) (models.Message, error)

	// Utils function that checks if user is partecipant in a conversation
	CheckUserConversation(user_id int64, conversation_id int64) (bool, error)

	// Get comment from a message
	GetComments(user_id int64, conversation_id int64, message_id int64) ([]models.Comment, error)

	// Create a comment under a message
	CreateComment(user_id int64, conversation_id int64, message_id int64, content string) (models.Comment, error)

	// Delete a comment
	DeleteComment(user_id int64, conversation_id int64, comment_id int64) error

	// Allows to add a user in a group chat
	AddGroup(user_id int64, username string, conversation_id int64) error

	// Remove User / or left from group
	RemoveGroup(user_id int64, conversation_id int64, member_id int64) error

	// Edit group photo
	EditPhoto(conversation_id int64, groupPhoto []byte) error

	// Edit group name
	EditName(conversation_id int64, groupName string) error

	// Edit profile name
	EditProfileName(user_id int64, username string) error

	// Edit profile photo
	EditProfilePhoto(user_id int64, photo []byte) error

	// Get Users by query
	GetUsers(names string) []models.User

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	TableMapping := map[string]string{
		"Users":         usersTableCreationStatement,
		"Conversations": conversationsTableCreationStatement,
		"Partecipants":  partecipantsTableCreationStatement,
		"Messages":      messagesTableCreationStatement,
		"Comments":      commentsTableCreationStatement,
	}

	for tableName, tableCreationStatement := range TableMapping {
		// Check if the table exist, otherwise we create it
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name= ? ;`, tableName).Scan(&tableName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err = db.Exec(tableCreationStatement)
				if err != nil {
					return nil, errors.New("error building table " + tableName)
				}
			} else {
				return nil, errors.New("error checking table " + tableName)
			}
		}
	}

	// query := `
	// 	INSERT INTO Users (username, profile_photo) VALUES
	// 	('user1', NULL),
	// 	('user2', NULL),
	// 	('user3', NULL),
	// 	('user4', NULL);

	// 	INSERT INTO Conversations (name, photo, type) VALUES
	// 	('Conversation1_User1', NULL, 'group'),
	// 	('Conversation2_User1', NULL, 'group'),
	// 	('Conversation1_User2', NULL, 'group'),
	// 	('Conversation2_User2', NULL, 'group');

	// 	INSERT INTO Partecipants (user_id, conversation_id) VALUES
	// 	(1, 1), (1, 2),
	// 	(2, 1), (2, 3),
	// 	(3, 2), (3, 3),
	// 	(4, 4);

	// 	INSERT INTO Messages (conversation_id, user_id, content, media, type, timestamp, status, isForwarded) VALUES
	// 	(1, 1, 'Hello from User1', NULL, 'text', 1691234567, 'sent', 0),
	// 	(1, 2, 'Hello from User2', NULL, 'text', 1691235567, 'sent', 0),
	// 	(1, 1, 'How are you?', NULL, 'text', 1691236567, 'sent', 0),
	// 	(1, 2, 'I am fine, thank you!', NULL, 'text', 1691237567, 'sent', 0),

	// 	(2, 1, 'Welcome to Conversation2', NULL, 'text', 1691238567, 'sent', 0),
	// 	(2, 3, 'Thanks!', NULL, 'text', 1691239567, 'sent', 0),
	// 	(2, 1, 'Let’s discuss our project.', NULL, 'text', 1691240567, 'sent', 0),
	// 	(2, 3, 'Sure!', NULL, 'text', 1691241567, 'sent', 0),

	// 	(3, 2, 'Conversation1_User2 message1', NULL, 'text', 1691242567, 'sent', 0),
	// 	(3, 3, 'Conversation1_User2 message2', NULL, 'text', 1691243567, 'sent', 0),
	// 	(3, 2, 'Conversation1_User2 message3', NULL, 'text', 1691244567, 'sent', 0),
	// 	(3, 3, 'Conversation1_User2 message4', NULL, 'text', 1691245567, 'sent', 0),

	// 	(4, 4, 'Hello from User4 in Conversation2_User2', NULL, 'text', 1691246567, 'sent', 0),
	// 	(4, 4, 'Second message User4', NULL, 'text', 1691247567, 'sent', 0),
	// 	(4, 4, 'Third message User4', NULL, 'text', 1691248567, 'sent', 0),
	// 	(4, 4, 'Fourth message User4', NULL, 'text', 1691249567, 'sent', 0);
	// 	`

	// _, err := db.Exec(query)
	// if err != nil {
	// 	log.Fatalf("Errore durante l'esecuzione della query: %v", err)
	// }

	// fmt.Println("Database popolato con successo!")

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
