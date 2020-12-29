DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;

\c test_db;

CREATE TABLE IF NOT EXISTS job(id serial PRIMARY KEY, company varchar(255) NOT NULL, title text NOT NULL, url text NOT NULL, dateadded varchar(255) NOT NULL);
CREATE TABLE IF NOT EXISTS user_list(id serial PRIMARY KEY, email varchar(255) NOT NULL, password varchar(255) NOT NULL);
CREATE TABLE IF NOT EXISTS user_job(id serial PRIMARY KEY, user_id int, job_id int);

ALTER TABLE user_job
ADD CONSTRAINT user_id_fk
FOREIGN KEY (user_id) REFERENCES user_list (id);

ALTER TABLE user_job
ADD CONSTRAINT job_id_fk
FOREIGN KEY (job_id) REFERENCES job (id);

INSERT INTO user_list(email,password) SELECT 'kenichiderumo@gmail.com', '2030';