-- 授权 root 用户可以远程链接
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'root';
flush privileges;

create database if not exists scheduling_system default character set utf8mb4 collate utf8mb4_general_ci;

use scheduling_system;
-- 用户表：存储系统用户信息
create table if not exists user (
                                    id int auto_increment primary key, -- 用户ID，自增主键
                                    username varchar(20) null, -- 用户名
                                    password varchar(80) default '666666' null, -- 密码，默认值为'666666'
                                    gender tinyint(1) default 0 null, -- 性别，0表示男性，1表示女性
                                    name varchar(20) default '匿名用户' null, -- 用户姓名，默认为'匿名用户'
                                    phone varchar(20) default '0' null, -- 手机号码，默认为'0'
                                    work_year int default 0 null, -- 工作年限，默认为0
                                    group_id int default 0 null, -- 用户分组ID，默认为0
                                    post_id int default 0 null, -- 职位ID，默认为0
                                    constraint user_username_uindex unique (username) -- 用户名唯一索引约束
);

-- 职位表：存储系统职位信息
create table if not exists post (
                                    id int auto_increment primary key, -- 职位ID，自增主键
                                    name varchar(20) null, -- 职位名称
                                    number int -- 职位编号
);

-- 手术室表：存储手术室信息
create table if not exists operating_room (
                                              id int auto_increment primary key, -- 手术室ID，自增主键
                                              name varchar(20) null, -- 手术室名称
                                              category int -- 手术室类别
);

-- 手术类别表：存储手术类别信息
create table if not exists category (
                                        id int auto_increment primary key, -- 类别ID，自增主键
                                        name varchar(20) null -- 类别名称
);

-- 手术表：存储手术信息
create table if not exists surgery (
                                       id int auto_increment primary key, -- 手术ID，自增主键
                                       name varchar(20) null, -- 手术名称
                                       category_id int, -- 手术类别ID，外键关联手术类别表
                                       room_id int, -- 手术室ID，外键关联手术室表
                                       start_time datetime, -- 手术开始时间
                                       end_time datetime, -- 手术结束时间
                                       state int -- 手术状态
);

-- 医护士手术关系表：存储医护人员与手术之间的关系
create table if not exists surgery_user (
                                            surgery_id bigint not null, -- 手术ID，外键
                                            user_id bigint not null, -- 用户ID，外键
                                            primary key (surgery_id, user_id), -- 主键由手术ID和用户ID组成
                                            constraint fk_surgery_user_surgery foreign key (surgery_id) references surgery (id), -- 外键约束，关联手术表的ID字段
                                            constraint fk_surgery_user_user foreign key (user_id) references user (id) -- 外键约束，关联用户表的ID字段
);
