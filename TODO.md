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
