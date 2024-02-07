# Design Notes

Original sketch of the concept is [here](https://github.com/benjaminjkraft/notes/blob/master/github-patch-workflow.md).

## When to destack

In theory the algorithm should be able to create non-stacked PRs from the patch-stack if there are no conflicts. But that produces added complexity because we may need to re-stack those PRs in the future. So for now let's keep everything stacked; auto-destacking can be in v0.2.

If we don't destack, other than the user-visible interactive rebase in `pgh pull`, we never need to manipulate trees, just branches and commits, because the trees always exactly match what's on GitHub. (It's possible this also forces us to say you have to pull before you push or something, but I don't even think so.) Once we destack, we do need to be able to manipulate trees, which will require nontrivial merges. This means we probably can't use `go-git` right now and instead have to use something with full `libgit2` or shell out.

## Operations we need

### pgh pr

1. We must be one commit ahead of the tip of our (possibly trivial) stack.
2. We push the existing stack, if any.
3. We cherry-pick the top commit onto the head of the top branch of the stack, push that to a new branch, and create a PR against the previously top branch of the stack for that new branch.

(This could be merged into push, but I think it will be a better UX without.)

### pgh push

1. We must be at the tip of our nontrivial stack.
2. Do the dance to convert rebases to merges at each point on the stack, and push each merge.
3. TODO: where do we get the commit messages? How do we tell for which commits you need them? (Maybe we want to say you shouldn't be at the tip of the stack? In theory merges are good about making it so we can update branches lazily. Rethink this whole command probably.)

Note: in principle there should be no need to pull before you push; if the bottom PR was merged we'll just skip its branch I guess?

### pgh pull

1. We must be at the tip of our nontrivial stack.
2. Handle any merged PRs by interactive-rebasing onto main and removing them from the stack. (This is where the user may need to resolve conflicts; we'll probably just shell out so they can do an ordinary rebase.)
3. In the future, this should also handle the case where someone else (e.g. "accept suggestion") pushed commits to the PRs; we can squash those in in our interactive rebase.

## Our metadata

We need to keep, for each commit in the stack:

- the corresponding PR (to handle when it was merged)
- the corresponding branch (or get this from the PR)
- the previous commit in the stack (to update correctly if you insert a commit in the middle?)
- the previous version of this commit (I think only needed for destacking?)

We need to put something in annotations so it is preserved when the user does a rebase. But putting all of this is probably too much. Probably the best thing is to put in a "pull request URL", and then just keep everything else indexed by that in git config or a file in `.git` or something?

## Open questions

- Do we need to do anything fancy to allow multiple stacks?
- Do we need to prevent running `pgh` when you're not at the tip? Or what do we need to do differently? (In fact we need to allow `pgh pr` when not at the tip, probably!)
- How do we name branches?
