CREATE TABLE time_slots
(
    "index" INTEGER PRIMARY KEY,
    startTime CHARACTER(5),
    endTime CHARACTER(5)
);

INSERT INTO time_slots("index", startTime, endTime) VALUES (1, '08:30', '10:05');
INSERT INTO time_slots("index", startTime, endTime) VALUES (2, '10:20', '11:55');
INSERT INTO time_slots("index", startTime, endTime) VALUES (3, '12:10', '13:45');
INSERT INTO time_slots("index", startTime, endTime) VALUES (4, '14:15', '15:50');
INSERT INTO time_slots("index", startTime, endTime) VALUES (5, '16:05', '17:40');
INSERT INTO time_slots("index", startTime, endTime) VALUES (6, '17:50', '19:25');
INSERT INTO time_slots("index", startTime, endTime) VALUES (7, '19:35', '21:10');


CREATE TABLE classes
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(100),
    "group" VARCHAR(10),
    subgroup INTEGER CHECK (subgroup BETWEEN 0 AND 4),
    building VARCHAR(20),
    "type" VARCHAR(5),
    room VARCHAR(20),
    professors VARCHAR(100),
    notes VARCHAR(100),
    regularity INTEGER CHECK (regularity BETWEEN 1 AND 3),
    "day" INTEGER CHECK ("day" BETWEEN 1 AND 6),

    "index" INTEGER REFERENCES time_slots("index")
);
