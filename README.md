# BallotCleaner

This was designed specifically for the ASOSU 2022 election where data was ingested from Qualtrics. It was designed to clean the data and create a CSV file that can be used to determine the winners of the election.

Namely it has four main steps to clean the data:
1. Check against registrar data to see if the voter is registered as an ASOSU student (graduate or undergraduate at the Corvallis Campus)
2. Check against previously submitted ballots to see if the voter has already submitted a ballot
3. Check to make sure each voter voted for the correct House of Representatives candidates (either grad or undergrad)
4. Select the winners for each incentive (daily, weekly and overall)

## Command Line Usage

`program <day>` - Processes the data for the given day.

`program <start_day> <end_day>` - Processes the data for the given range of days (inclusive).

`./scripts/run.sh` will run the program and accept the above arguments

## Input

This program expects a folder named `data` containing the following files:
- `seed.txt` - a single line of text containing the seed for the random number generator used to select the winners
- `validVoters.csv` - a CSV file containing the valid voters in the form of `FIRST_NAME LAST_NAME	OSU_EMAIL	ONID_ID	G_UG_STATUS` (separated by tabs)
- `ballots` - a folder containing all the ballots submitted by the voters in the form of `<days_since_epoch>-whatever.csv`. It is expected that the files contain data for the day listed as well as all days prior. The format is too long to document here and must be customized for each election/ballot.
- `alreadyVoted` - a folder containing files in the form of `whatever-<days_since_epoch>-whatever.csv` which lists all the voters who have already voted on a given day. This data is deduped so there is no harm in having overlapping data. One ONID per line

## Output

Each step along the way, this program will output the data about what it did. The first number in each file is the step it corresponds to. The next text corresponds to the type of the data and then the final two numbers represent the start date and end date (inclusive) of the data. For example, `1-invalid-3-5.csv` represents the invalid data from step one for days 3, 4, and 5.

Each step also outputs a summary with the number of valid, invalid, and total votes processed, as well as any additional log information that might be useful.

Step 2 outputs an additional file that can be copied directly into the `alreadyVoted` folder of the input. Step 4 outputs an additional file proving the ONID IDs of the winners.
