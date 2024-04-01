DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS groups CASCADE;
DROP TABLE IF EXISTS appliances CASCADE;

CREATE TABLE users (
  id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text      NOT NULL UNIQUE,
  nickname text NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE groups (
  id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text      NOT NULL,
  user_id uuid NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  CONSTRAINT fk_user
    FOREIGN KEY(user_id) 
    REFERENCES users(id)
);

CREATE TABLE appliances (
  id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text      NOT NULL,
  wattage decimal NOT NULL,
  group_id uuid,
  user_id uuid,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  CONSTRAINT fk_group
    FOREIGN KEY(group_id) 
    REFERENCES groups(id),
  CONSTRAINT fk_user
    FOREIGN KEY(user_id) 
    REFERENCES users(id)
);
