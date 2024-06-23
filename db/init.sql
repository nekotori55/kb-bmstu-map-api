CREATE DATABASE kb_bmstu_map_api;

CREATE TABLE time_slots
(
    id SERIAL PRIMARY KEY,
    startTime TIME,
    endTime TIME
);

INSERT INTO time_slots(startTime, endTime) VALUES ('08:30', '10:05');
INSERT INTO time_slots(startTime, endTime) VALUES ('10:20', '11:55');
INSERT INTO time_slots(startTime, endTime) VALUES ('12:10', '13:45');
INSERT INTO time_slots(startTime, endTime) VALUES ('14:15', '15:50');
INSERT INTO time_slots(startTime, endTime) VALUES ('16:05', '17:40');
INSERT INTO time_slots(startTime, endTime) VALUES ('17:50', '19:25');
INSERT INTO time_slots(startTime, endTime) VALUES ('19:35', '21:10');


CREATE TABLE schedule
(
    id SERIAL PRIMARY KEY,
    title CHARACTER(100),
    "group" INTEGER CHECK ("group" BETWEEN 1 AND 4),
    building CHARACTER(20),
    room CHARACTER(20),
    professors CHARACTER(100),
    notes CHARACTER(100),
    regularity INTEGER CHECK (regularity BETWEEN 1 AND 3),
    day_number INTEGER CHECK (day_number BETWEEN 1 AND 6),

    time_slot INTEGER REFERENCES time_slots(id)
)
