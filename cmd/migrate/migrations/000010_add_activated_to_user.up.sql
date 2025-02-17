alter TABLE
    users
ADD
COLUMN is_active boolean not NULL DEFAULT false;