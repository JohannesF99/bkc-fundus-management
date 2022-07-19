CREATE TABLE entries (
	id          INT         NOT NULL    AUTO_INCREMENT,
	member_id   INT         NOT NULL,
	item_id     INT         NOT NULL,
	capacity    INT         NOT NULL    CHECK (capacity>0),
    created 	TIMESTAMP 	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
    modified 	TIMESTAMP 	NOT NULL 	DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (member_id) REFERENCES members(id),
	FOREIGN KEY (item_id) REFERENCES items(id),
	PRIMARY KEY (id)
);