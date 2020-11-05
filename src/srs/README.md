This is nothing special - load the student data from a csv (omitted in the repo for privacy reasons), seed the PRNG with the current time, get a random number and check 
if it's been used before using a `HashMap` because I guess Go doesn't have `Set`s, and keep going until the sample has been picked. This is not the most time effecient approach,
but I didn't want to bother with shuffling or anything too crazy, so I went with this. It'll work well so long as the sample doesn't get too large. Then obviously output it into
separate csv files for the treatments.
