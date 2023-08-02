DROP TABLE IF EXISTS category;
CREATE TABLE category (
  id		INT AUTO_INCREMENT NOT NULL,
  name		VARCHAR(127) NOT NULL UNIQUE,
  description	VARCHAR(255),
  PRIMARY KEY (`id`)
);

DROP TABLE IF EXISTS tag;
CREATE TABLE tag (
  id	INT AUTO_INCREMENT NOT NULL,
  name	VARCHAR(127) NOT NULL,
  PRIMARY KEY (`id`)
);

DROP TABLE IF EXISTS expense;
CREATE TABLE expense (
  id		INT AUTO_INCREMENT NOT NULL,
  category_id	INT NOT NULL,
  amount	DECIMAL(10,2) NOT NULL,
  date		DATE NOT NULL,
  location	VARCHAR(255),
  notes		VARCHAR(255),
  PRIMARY KEY (`id`),
  FOREIGN KEY (`category_id`) REFERENCES category(`id`)
);

DROP TABLE IF EXISTS expense_tags;
CREATE TABLE expense_tags (
  id		INT AUTO_INCREMENT NOT NULL,
  expense_id	INT NOT NULL,
  tag_id	INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`expense_id`) REFERENCES expense(`id`),
  FOREIGN KEY (`tag_id`) REFERENCES tag(`id`)
);
