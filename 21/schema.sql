CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

create table comments (
    id serial primary key,
    task_id int not null references tasks(id) on delete cascade,
    author int not null references users(id) on delete cascade,
    text text not null,
    created_at timestamp default now()
);