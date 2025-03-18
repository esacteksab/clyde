> There are two hard problems in computer science: cache invalidation, naming things and off-by-one errors.
> -- [Leon Bambrick](https://twitter.com/secretGeek/status/7269997868)

# Clyde

Right Turn Clyde.

Another experiment. My very naive attempt at detecting "hijacked" or "squatted" Go modules. Still a work in progress. Uses two metrics, the age of the repository and whether or not it's a fork. There is likely more nuance to suss out. Maybe number of stars, number of commits after it was forked. Maybe traversing the fork(s) to identify patterns? Probably rub some AI on it.

I may take this a little further. I may never revisit this. But it's an itch I needed to scratch. Only time will tell.

## References

For greater insight or understanding, I stumbled upon these two blog posts:

- [https://mhouge.dk/blog/rogue-one-a-malware-story/](https://mhouge.dk/blog/rogue-one-a-malware-story/)
- [https://alexandear.github.io/posts/2025-02-28-malicious-go-programs/](https://alexandear.github.io/posts/2025-02-28-malicious-go-programs/)

I originally heard about this on Reddit in [r/golang](https://www.reddit.com/r/golang/comments/1jbzuot/someone_copied_our_github_project_made_it_look/)

Becareful out there!
