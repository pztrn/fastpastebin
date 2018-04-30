package migrations

import (
	// stdlib
	"database/sql"
)

func InitialUp(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE TABLE `pastes` (`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Paste ID', `title` text NOT NULL COMMENT 'Paste title', `data` longtext NOT NULL COMMENT 'Paste data', `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Paste creation timestamp', `keep_for` int(4) NOT NULL DEFAULT 1 COMMENT 'Keep for integer. 0 - forever.', `keep_for_unit_type` int(1) NOT NULL DEFAULT 1 COMMENT 'Keep for unit type. 1 - minutes, 2 - hours, 3 - days, 4 - months.', PRIMARY KEY (`id`), UNIQUE KEY `id` (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='Pastes';")
	if err != nil {
		return err
	}

	return nil
}
