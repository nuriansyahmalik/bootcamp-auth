CREATE TABLE IF EXISTS user (
    id varchar(36) PRIMARY KEY,
    username varchar(255),
    name varchar(255),
    role varchar(255),
    password
    );

CREATE TABLE  IF EXISTS course (
    id varchar(36) PRIMARY KEY,
    title varchar(255),
    content varchar(255),
    user_id varchar(255)
    );

create index idx_id on user(id);
alter table course add foreign key (user_id) references user(id)
