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
- then automate the "walk-up-the-stack" part
- then actually tie it together, do all the silly mechanics

for fake-merge:

for walk-up:
- worktree? prob easier to not, you want to end up at the top anyway, be able to resolve conflicts, etc. (maybe option later.)
- how does pause/resume for conflicts work? (for stack: just run again, tree would be nice though. add later.)
- do rebase for unpushed branches? eh, maybe not, can always do that on top
- do auto-push? as option maybe to start, sniffing will be fragile
