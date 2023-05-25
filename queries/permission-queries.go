package queries

import (
	"database/sql"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go.uber.org/zap"
)

func checkGrantedUserExists(logger *zap.Logger, db *sql.DB, userID int64, grantedUser *tgbotapi.User) (bool, error) {
	kGetGrantedUser := `
	SELECT COUNT(*)
	FROM permissions
	WHERE user_id = ?
	  AND granted_user_id = ?
	`

	rows, err := db.Query(kGetGrantedUser, userID, grantedUser.ID)
	if err != nil {
		logger.Error("Error querying kGetGrantedUser", zap.Error(err))
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var cnt int
		err := rows.Scan(&cnt)
		if err != nil {
			logger.Error("Error scanning kGetGrantedUser", zap.Error(err))
			return false, err
		}
		logger.Debug("Fetched from permissions. cnt = " + fmt.Sprint(cnt))
		return cnt > 0, nil
	}
	logger.Debug("Could not get rows.Next() while fetching permissions")
	return false, err
}

func insertGrantedUser(logger *zap.Logger, db *sql.DB, user, grantedUser *tgbotapi.User) error {
	InsertUser(logger, db, user)
	InsertUser(logger, db, grantedUser)

	kInsertNewGrantedUser := `
	INSERT INTO permissions (user_id, granted_user_id)
	VALUES (?, ?);
	`

	_, err := db.Exec(kInsertNewGrantedUser, user.ID, grantedUser.ID)
	if err != nil {
		logger.Error("Error inserting note into database", zap.Error(err))
		return err
	}

	logger.Debug(fmt.Sprintf("Successfully inserted user_id=%d, granted_user_id=%d", user.ID, grantedUser.ID))
	return nil
}

func GetUserPermissions(logger *zap.Logger, db *sql.DB, userID int64) ([]string, error) {
	kGetPermissions := `
	SELECT u.name
	FROM permissions p
	INNER JOIN users u ON p.granted_user_id = u.id
	WHERE p.user_id = ?;
	`

	rows, err := db.Query(kGetPermissions, userID)
	if err != nil {
		logger.Error("Error while getting user permissions", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var granted_users []string
	for rows.Next() {
		var login string
		err := rows.Scan(&login)
		if err != nil {
			logger.Error("Error while scanning row", zap.Error(err))
			return nil, err
		}

		granted_users = append(granted_users, login)
	}

	logger.Info("Granted users: [" + strings.Join(granted_users, ", ") + "]")
	return granted_users, nil
}

func CheckUserGavePermission(logger *zap.Logger, db *sql.DB, userID int64, grantedUser string) (bool, int64, error) {
	kHasPermission := `
	SELECT p.granted_user_id, p.user_id
	FROM permissions p
	INNER JOIN users u ON u.id = p.user_id
	WHERE u.name = ?
	`

	rows, err := db.Query(kHasPermission, grantedUser)
	if err != nil {
		logger.Error("Error while getting user permission", zap.Error(err))
		return false, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var granted_user_id, user_id int64
		err := rows.Scan(&granted_user_id, &user_id)
		if err != nil {
			logger.Error("Error while scanning row", zap.Error(err))
			return false, 0, err
		}
		return granted_user_id == userID, user_id, nil
	}

	// No rows - user did not grant an access
	return false, 0, nil
}

func ChaeckAndInsertNewGrantedUser(logger *zap.Logger, db *sql.DB, user, grantedUser *tgbotapi.User) (bool, error) {
	exists, err := checkGrantedUserExists(logger, db, user.ID, grantedUser)
	if err != nil {
		logger.Error("Check if granted user exists failed", zap.Error(err))
		return false, err
	}

	if exists {
		logger.Info("The user is already granted")
		return false, err
	}

	err = insertGrantedUser(logger, db, user, grantedUser)
	if err != nil {
		logger.Error("Insertion of granted user failed", zap.Error(err))
		return false, err
	}

	return true, nil
}

func DeleteAllGrantedUsers(logger *zap.Logger, db *sql.DB, userID int64) error {
	kDeleteGrantedUsers := `
	DELETE FROM permissions
	WHERE user_id = ?
	`

	_, err := db.Exec(kDeleteGrantedUsers, userID)
	if err != nil {
		logger.Error("Error inserting note into database", zap.Error(err))
	}
	return err
}
