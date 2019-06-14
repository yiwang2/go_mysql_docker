CREATE DATABASE IF NOT EXISTS poppulodb;
USE poppulodb;
CREATE TABLE IF NOT EXISTS tickets (id BIGINT(20) NOT NULL PRIMARY KEY,status varchar(40) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=latin1;
CREATE TABLE IF NOT EXISTS ticket_line (ticket_id BIGINT(20) NOT NULL, line varchar(400) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=latin1;

insert into tickets (id, status) values (1555759165502, "UNCHECKED");
insert into ticket_line (ticket_id, line) values(1555759165502 ,'{\"values\":[2,2,2],\"result\":5}');