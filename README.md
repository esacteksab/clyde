> There are two hard problems in computer science: cache invalidation, naming things and off-by-one errors. <!--markdownlint-disable MD041 -->
>
> -- [Leon Bambrick](https://twitter.com/secretGeek/status/7269997868)

# Clyde

Right Turn Clyde.

Another experiment. My very naive attempt at detecting "hijacked" or "squatted" Go modules. Still a work in progress. Uses two metrics, the age of the repository and whether or not it's a fork. There is likely more nuance to suss out. Maybe number of stars, number of commits after it was forked. Maybe traversing the fork(s) to identify patterns? Probably rub some AI on it.

I may take this a little further. I may never revisit this. But it's an itch I needed to scratch. Only time will tell.

## March 18, 2025 Now with a gratuitous amount of emojis

Cause apparently I turned in a teenage girl. Oh, and caching! You're welcome GitHub! (_p.s sorry about the rate limit thing, hope we can still be friends._)

A Happy Module

```text
====================
🔑 Found GITHUB_TOKEN in environment.
🔐 Using authenticated GitHub client
✅ RESPONSE SERVED FROM CACHE
⚠️ No Authorization header found!
Rate limit: 5000/5000, resets at 2025-03-18-23:10:56
✅ Using authenticated rate limits (5000+/hour)
Module is 3883 days old.
Module is: github.com/google/btree
Repo was created at: 2014-07-31
🍰 Repo is not a fork
Repo last updated at: 2025-03-19
✨ Module has a score of: 0.000000 out of 100.
```

A Less Happy Module

```text
====================
🔑 Found GITHUB_TOKEN in environment.
🔐 Using authenticated GitHub client
✅ RESPONSE SERVED FROM CACHE
⚠️ No Authorization header found!
Rate limit: 5000/5000, resets at 2025-03-18-23:10:55
✅ Using authenticated rate limits (5000+/hour)
Module is 23 days old.
Module is: github.com/esacteksab/sausage-factory
Repo was created at: 2025-02-23
🍰 Repo is not a fork
Repo last updated at: 2025-03-16
✨ Module has a score of: 38.330000 out of 100.
```

A Sad Module

```text
====================
🔑 Found GITHUB_TOKEN in environment.
🔐 Using authenticated GitHub client
✅ RESPONSE SERVED FROM CACHE
⚠️ No Authorization header found!
Rate limit: 5000/5000, resets at 2025-03-18-23:10:55
✅ Using authenticated rate limits (5000+/hour)
Module is 3694 days old.
Module is: github.com/launchdarkly/httpcache
Repo was created at: 2015-02-06
🍴 Repo is a fork
Repo last updated at: 2024-12-09
💩 Module has a score of: 50.000000 out of 100.
```

A More Sad Module

```text
====================
🔑 Found GITHUB_TOKEN in environment.
🔐 Using authenticated GitHub client
✅ RESPONSE SERVED FROM CACHE
⚠️ No Authorization header found!
Rate limit: 5000/5000, resets at 2025-03-18-23:10:56
✅ Using authenticated rate limits (5000+/hour)
Module is 1 days old.
Module is: github.com/esacteksab/sshalert
Repo was created at: 2025-03-18
🍴 Repo is a fork
Repo last updated at: 2025-03-18
⛔ Module has a score of: 101.670000 out of 100.
```

Clearly I have some formatting to sort out with `float64`. This alone should further support not using this. Look who created it.

Caching currently writes to current working directory. Will likely put it in [os.UserCacheDir](https://pkg.go.dev/os#UserCacheDir) I think. It existing there, being more global would allow the utility of this to be a shared resource across repos/projects rather than caching being isolated to the current project.

## References

For greater insight or understanding, I stumbled upon these two blog posts:

- [https://mhouge.dk/blog/rogue-one-a-malware-story/](https://mhouge.dk/blog/rogue-one-a-malware-story/)
- [https://alexandear.github.io/posts/2025-02-28-malicious-go-programs/](https://alexandear.github.io/posts/2025-02-28-malicious-go-programs/)

I originally heard about this on Reddit in [r/golang](https://www.reddit.com/r/golang/comments/1jbzuot/someone_copied_our_github_project_made_it_look/)

Becareful out there!

### Attribution

A shoutout to my buddy Claude! They were often wrong the first few times, but we persevered and came through the otherside. Some Go packages have non-existent documentation like I'm supposed to know this shit or something. Damn, help a n00b out, throw us a bone with a quick how-to, even if it is littered with shitty poorly named variables that just confuse the reader further and don't reflect a real-world use case. _Anything_ is better than nothing.
