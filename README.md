> There are two hard problems in computer science: cache invalidation, naming things and off-by-one errors. <!--markdownlint-disable MD041 -->
>
> -- [Leon Bambrick](https://twitter.com/secretGeek/status/7269997868)

# Clyde

_Right Turn Clyde!_

---

Another experiment. My very naive attempt at detecting "hijacked" or "squatted" Go modules. Still a work in progress. Uses two metrics, the age of the repository and whether it's a fork. There is likely more nuance to suss out. Maybe number of stars, number of commits after it was forked. Maybe traversing the fork(s) to identify patterns? Probably rub some AI on it.

I may take this a little further. I may never revisit this. But it's an itch I needed to scratch. Only time will tell.

### A little more scratching

A Happy Module

```text
🔑 Found GITHUB_TOKEN in environment.
🔐 Using authenticated GitHub client
✅ RESPONSE SERVED FROM CACHE
⚠️ No Authorization header found!
Rate limit: 4933/5000, resets at 2025-03-21-20:40:24
✅ Using authenticated rate limits (5000+/hour)
Module is: github.com/google/go-github
🍰 Repo is not a fork
Repo was created at: 2013-05-24
Repo last updated at: 2025-03-21
Module is 4319 days old.
✨ Module has a score of: 0.00 out of 100.

B===================================================D
```

A Less Happy Module

```text
🔑 Found GITHUB_TOKEN in environment.
🔐 Using authenticated GitHub client
✅ RESPONSE SERVED FROM CACHE
⚠️ No Authorization header found!
Rate limit: 4933/5000, resets at 2025-03-21-20:40:24
✅ Using authenticated rate limits (5000+/hour)
Module is: github.com/launchdarkly/httpcache
🍴 Repo is a fork
Repo was created at: 2015-02-06
Repo last updated at: 2024-12-09
Module is 3697 days old.
💩 Module has a score of: 50.00 out of 100.

B===================================================D
```

A Sad Module

```text
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

B===================================================D
```

A More Sad Module

```text
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

B===================================================D
```

A Module Not Hosted on GitHub

```text

❓ Module golang.org/x/text not hosted on GitHub.
B===================================================D
```

Clearly I have some formatting to sort out with `float64`. This alone should further support not using this. Look who created it.

See the [CHANGELOG](./CHANGELOG.md) for additional details.

## References

For greater insight or understanding, I stumbled upon these two blog posts:

- [https://mhouge.dk/blog/rogue-one-a-malware-story/](https://mhouge.dk/blog/rogue-one-a-malware-story/)
- [https://alexandear.github.io/posts/2025-02-28-malicious-go-programs/](https://alexandear.github.io/posts/2025-02-28-malicious-go-programs/)

I originally heard about this on Reddit in [r/golang](https://www.reddit.com/r/golang/comments/1jbzuot/someone_copied_our_github_project_made_it_look/)

Be careful out there!

### Attribution

A shout-out to my buddy Claude! They were often wrong the first few times, but we persevered and came through the other side. Some Go packages have non-existent documentation like I'm supposed to know this shit or something. Damn, help a n00b out, throw us a bone with a quick how-to, even if it is littered with shitty poorly named variables that just confuse the reader further and don't reflect a real-world use case. _Anything_ is better than nothing. Claude's largest contribution was providing an example of using `httpcache` with `go-github`.
