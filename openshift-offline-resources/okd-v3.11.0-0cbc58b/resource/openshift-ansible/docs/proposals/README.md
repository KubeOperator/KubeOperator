# OpenShift-Ansible Proposal Process

## Proposal Decision Tree
TODO: Add details about when a proposal is or is not required. 

## Proposal Process
The following process should be followed when a proposal is needed:

1. Create a pull request with the initial proposal
  * Use the [proposal template][template]
  * Name the proposal using two or three topic words with underscores as a separator (i.e. proposal_template.md)
  * Place the proposal in the docs/proposals directory
2. Notify the development team of the proposal and request feedback
3. Review the proposal on the OpenShift-Ansible Architecture Meeting
4. Update the proposal as needed and ask for feedback
5. Approved/Closed Phase
  * If 75% or more of the active development team give the proposal a :+1: it is Approved
  * If 50% or more of the active development team disagrees with the proposal it is Closed
  * If the person proposing the proposal no longer wishes to continue they can request it to be Closed
  * If there is no activity on a proposal, the active development team may Close the proposal at their discretion
  * If none of the above is met the cycle can continue to Step 4.
6. For approved proposals, the current development lead(s) will:
  * Update the Pull Request with the result and merge the proposal
  * Create a card on the Cluster Lifecycle [Trello board][trello] so it may be scheduled for implementation.

[template]: proposal_template.md
[trello]: https://trello.com/b/wJYDst6C
