immediate TODO:
- test harness
  - can I start with in-memory repo?
  - ops on local repo
  - fake GH
  - real GH
- diff command to put up first PR
  - for github base is a branch; need to figure out right semantics but maybe
    start with just "whatever branch points at HEAD^", which is totally wrong
    but probably 80% good enough to start with
  - git push
  - make PR via API
- commit -> PR # association
  - use annotations?
  - then just filter out of PR messages
- will not having merge implemented be an issue?
  - prob not actually, just need to create merge commits but
    we know what all the trees should look like I think?

thoughts on what to prioritize Jan 2024:
- make it useful today
- start with a command to do the `merge -X` thing, since that's not really possible rn
- then automate the "walk-up-the-stack" part (without worktree this is not that hard)
- then actually tie it together, do all the silly mechanics

for fake-merge:
- kinda confusing if upstream isn't main (e.g. you're on -2 and forgot to update it)

for merge-up:
- option to push all modified branches
