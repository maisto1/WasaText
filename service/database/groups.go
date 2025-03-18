package database

import "errors"

func (db *appdbimpl) AddGroup(user_id int64, username string, conversation_id int64) error {
	var isGroup bool
	user_partecipant := db.GetUsers(username)
	if len(user_partecipant) == 0 {
		return errors.New("participant not found")
	}
	partecipant_id := user_partecipant[0].User_id
	isValid, err := db.CheckPrivateConversation(user_id, partecipant_id)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("user doesn't have private conversation")
	}
	err = db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1
            FROM Conversations
            WHERE conversation_id = ? AND conversation_type = 'group'
        )
    `, conversation_id).Scan(&isGroup)
	if err != nil {
		return err
	}
	if !isGroup {
		return errors.New("this isn't a group chat")
	}
	_, err = db.c.Exec("INSERT INTO Partecipants (user_id, conversation_id) VALUES (?,?)", partecipant_id, conversation_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveGroup(user_id int64, conversation_id int64, member_id int64) error {
	var isGroup bool
	err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1
            FROM Conversations
            WHERE conversation_id = ? AND conversation_type = 'group'
        )
    `, conversation_id).Scan(&isGroup)
	if err != nil {
		return err
	}
	if !isGroup {
		return errors.New("this isn't a group chat")
	}
	_, err = db.c.Exec("DELETE FROM Partecipants WHERE user_id = ? AND conversation_id = ?;", member_id, conversation_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) EditName(conversation_id int64, groupName string) error {
	var isGroup bool
	err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1
            FROM Conversations
            WHERE conversation_id = ? AND conversation_type = 'group'
        )
    `, conversation_id).Scan(&isGroup)
	if err != nil {
		return err
	}
	if !isGroup {
		return errors.New("this isn't a group chat")
	}
	_, err = db.c.Exec("UPDATE Conversations SET name = ? WHERE conversation_id = ?", groupName, conversation_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) EditPhoto(conversation_id int64, groupPhoto []byte) error {
	var isGroup bool
	err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1
            FROM Conversations
            WHERE conversation_id = ? AND conversation_type = 'group'
        )
    `, conversation_id).Scan(&isGroup)
	if err != nil {
		return err
	}
	if !isGroup {
		return errors.New("this isn't a group chat")
	}
	_, err = db.c.Exec("UPDATE Conversations SET conversation_photo = ? WHERE conversation_id = ?", groupPhoto, conversation_id)
	if err != nil {
		return err
	}
	return nil
}
