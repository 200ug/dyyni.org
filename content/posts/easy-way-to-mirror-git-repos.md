---
title: "An easy way to mirror git repositories"
date: "2026-03-09T19:00:56+02:00"
draft: false
post_number: "006"
---

Due to the (relatively recent) enshittification of big tech platforms, most notably Microsoft (but also, unfortunately, platforms like YouTube), I've decided to move to primarily using Codeberg instead of GitHub to store my repositories. I'm also planning to host my own cgit instance once I get my home server put together, but the recent upward development in storage device prices has seriously hindered the progress on that front. Thus Codeberg will be a good stepping stone for this purpose, being a nonprofit and European.

Codeberg and other Forgejo-driven platforms have pretty great functionality for migrating repositories from other platforms, but if you'd also like to mirror future commits between these two repositories, the steps you need to take become more cumbersome. Basically, assuming you'd use the Codeberg repository as your primary remote from now on, you'd have to add the SSH clone URL of the GitHub repository, and copy the auto-generated (repository specific) SSH authentication keys from Codeberg to GitHub. For the initial migration step, this isn't that big of a deal (unless you have hundreds of repositories to migrate) as you'll be only doing that once, but I immediately felt that I wouldn't want to be doing this for every single new repository I create and want to mirror to both platforms. So, I decided to write a little shell script ["ferne"](https://codeberg.org/2ug/shellscripts/src/branch/master/ferne) to work around this problem.

![Is it worth the time? (xkcd)](/images/posts/easy-way-to-mirror-git-repos/is_it_worth_the_time.png)

Even though writing the script didn't take a long time, this might be one of those tasks not worth automating according to xkcd :^) Despite this, I really like the end result. 

Git as a protocol would support the so-called push-to-create functionality, which would allow repository creation from the command line with `git push`, but unfortunately neither Codeberg nor GitHub allows this ([to prevent using the platform as storage](https://codeberg.org/Codeberg/Community/issues/1643#issuecomment-2366256)). Luckily, this can be worked around using their REST APIs (requires storing the access tokens locally, but that shouldn't be an issue as long as you encrypt them with e.g. PGP like I did) to essentially achieve the following workflow:

1. List the mirror platforms and account usernames to a config file.
2. Get API tokens from the aforementioned platforms, put them into another config file, and encrypt it with your PGP key.
3. Run `ferne init` or `ferne private` to create a public or a private mirrored repository, respectively.

In the background, the script will check for naming conflicts, create new repositories on each mirrored platform, and configure the push URLs for local git (by default setting the first mirror as the primary remote). Notably the first two steps only need to be done during the tool setup, after which you can simply run the script during new repository initialization, and use `git push` to automagically push the pending commits to all mirrors without a hassle.

As a side note, nothing forces you to stop at just *one* mirror, you can go as deep as you wish.

![Why stop at just one mirror?](/images/posts/easy-way-to-mirror-git-repos/multimirror.png)

