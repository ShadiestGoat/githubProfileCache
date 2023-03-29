# Code Language Aggregator

A bit of a silly billy idea, but basically this project periodically fetches all the repos in your github profile, and counts up language use\*

You can then query for the aggregate results. 

\* - I added some language normalization tools, which aren't very accurate but are 'good enough'. These apply only to TS, JS, and Svelte, since they are bigger languages. I also distinguish between tsx & jsx, but that is a lot less accurate.

## Limitations

- It re-fetches only every 2 hours (cache)
- It doesn't count languages in:
  - gists
  - PRs
  - Forks
  - More?
  - Org projects
- Only 1 user
- Untested on Orgs

## Config

A `.env` file with 2 keys - `AUTHENTICATION`, which is where you place your github key (ghp_\*) and `PASSWORD`, which is where you write a password to authenticate requests. This can be anything!

Finally, this will be hosted on env variable `PORT`.
