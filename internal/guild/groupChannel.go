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

	_, err = db.Exec("UPDATE guilds SET group_channel = $1 WHERE id = $2", channelID, guildID)

	return
}
