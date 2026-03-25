CREATE TABLE subscription
(
    id serial unique not null,
    user_id uuid not null,
    service_name varchar(255) not null,
    price bigint not null,
    start_date varchar(255) not null,
    end_date date
);