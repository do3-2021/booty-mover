package guild

import "database/sql"

// Get group configuration channel for a guild
func GetGroupChannel(db *sql.DB, guildID string) (channelID string, err error) {

	row := db.QueryRow("SELECT group_channel FROM guilds WHERE id = $1", guildID)

	err = row.Scan(&channelID)

	return
}

// Set the group configuration channel for a guild
func SetGroupChannel(db *sql.DB, channelID, guildID string) (err error) {

	_, err = db.Exec(`INSERT INTO guilds (id, group_channel) 
	VALUES (?, ?)
	ON CONFLICT (id) DO UPDATE 
		SET group_channel = excluded.group_channel;`, channelID, guildID)

	return
}
