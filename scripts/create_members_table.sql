CREATE TABLE members(
	id 						INT 		NOT NULL 	AUTO_INCREMENT,
	name					varchar(255)NOT NULL,
	borrowed_item_count		INT 		NOT NULL 	DEFAULT 0,
	comment					TEXT		NOT NULL	DEFAULT '-',
	active                  BOOL        NOT NULL    DEFAULT TRUE,
	created 				TIMESTAMP 	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
	modified 				TIMESTAMP 	NOT NULL 	DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT borrowed_item_count_not_negativ CHECK (borrowed_item_count>=0),
	PRIMARY KEY (id)
);