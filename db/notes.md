Pet Daisy's database currently stores 37,023 users. Of those 37,023, only 36 have more than 100 pets. The database needs to be reworked to only store users who have a significant pet count (100 seems good for now.)

Users who aren't on the leaderboard also don't need their pets saved every pet. They can rely on autosave and save-on-disconnect


update may 7:
I keep hearing about bulk queries being better than for loops. I should look into this properly and will. soon (tm)

update june 30:

redis will solve the sqlite query issue; leaderboard deltas should help the data transfer issue.