
-- Create the Streams Table
CREATE TABLE IF NOT EXISTS streams(
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    updated_at    timestamp(3) default CURRENT_TIMESTAMP(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3) 
);


-- Create the Questions Table
CREATE TABLE IF NOT EXISTS questions(
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_type  ENUM( 'standard', 'poll')  DEFAULT 'standard', 
    question JSON NOT NULL, -- Question including the question text and answers are here
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    updated_at    timestamp(3) default CURRENT_TIMESTAMP(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3)
);

-- Create link table between questions and streams (im assuming a question can be associated to many streams)
CREATE TABLE IF NOT EXISTS question_stream(
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_id int NOT NULL,
    stream_id int NOT NULL,
    FOREIGN KEY (question_id)
        REFERENCES questions(id)
        ON DELETE CASCADE,
    INDEX stream_id (stream_id),
    FOREIGN KEY (stream_id)
        REFERENCES streams(id)    
);