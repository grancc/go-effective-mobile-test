ALTER TABLE subscription
    ALTER COLUMN start_date TYPE date USING to_date(start_date, 'MM-YYYY');
