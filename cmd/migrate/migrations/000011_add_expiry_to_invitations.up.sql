alter TABLE
user_invitations
ADD
COLUMN expiry TIMESTAMP(0) with time zone not NULL;