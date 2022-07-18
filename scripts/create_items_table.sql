CREATE TABLE items(
	id 						INT 		NOT NULL 	AUTO_INCREMENT,
	name					varchar(255)NOT NULL,
	capacity				INT 		NOT NULL	CHECK (capacity>0),
	available				INT 		NOT NULL,
	description				TEXT		NOT NULL	DEFAULT '-',
	created 				TIMESTAMP 	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
	modified 				TIMESTAMP 	NOT NULL 	DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT availability_less_or_equal_capacity CHECK (available<=capacity),
	PRIMARY KEY (id)
);