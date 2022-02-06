
CREATE TABLE users (
  idx SERIAL,
  id char(128) PRIMARY KEY,
  username CHAR(64),
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  email CHAR(256)
);

CREATE TABLE notes (
  idx SERIAL,
  id char(128) PRIMARY KEY,
  data TEXT NOT NULL,
  sort_order real NOT NULL DEfAULT 0,
  owner_id CHAR(128) NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);


ALTER TABLE notes ADD CONSTRAINT fk_notes_owner FOREIGN KEY (owner_id) REFERENCES users (id);
