ALTER TABLE subscription
    ALTER COLUMN start_date TYPE varchar(255) USING to_char(start_date, 'MM-YYYY');
