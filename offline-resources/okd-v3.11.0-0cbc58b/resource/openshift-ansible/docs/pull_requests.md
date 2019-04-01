# Pull Request process

Pull Requests in the `openshift-ansible` project follow a
[Continuous](https://en.wikipedia.org/wiki/Continuous_integration)
[Integration](https://martinfowler.com/articles/continuousIntegration.html)
process that is similar to the process observed in other repositories such as
[`origin`](https://github.com/openshift/origin).

Whenever a
[Pull Request is opened](../CONTRIBUTING.md#submitting-contributions), all
automated test jobs must be successfully run before the PR can be merged.

Some of these jobs are automatically triggered, e.g., Travis, PAPR, and
Coveralls. Other jobs need to be manually triggered by a member of the
[Team OpenShift Ansible Contributors](https://github.com/orgs/openshift/teams/team-openshift-ansible-contributors).

## Triggering tests

Members of the [Team OpenShift Ansible
Contributors](https://github.com/orgs/openshift/teams/team-openshift-ansible-contributors)
can trigger test jobs by adding a comment containing
`/ok-to-test`. For a full list of bot commands refer to the [Bot Command
Help](https://deck-ci.svc.ci.openshift.org/command-help?repo=openshift%2Fopenshift-ansible).

### Fedora tests

There are a set of tests that run on Fedora infrastructure. They are started
automatically with every pull request.

They are implemented using the [`PAPR` framework](https://github.com/projectatomic/papr).

To re-run tests, write a comment containing only `bot, retest this please`.

## Triggering merge

After a PR is properly reviewed and all test are passing, it can be
tagged for merge by a member of the [Team OpenShift Ansible
Contributors](https://github.com/orgs/openshift/teams/team-openshift-ansible-contributors)
by writing a comment containing `/lgtm` (looks good to me) anywhere in
the comment body.

Tagging a Pull Request with `/lgtm` puts it in an automated merge
queue. The
[@openshift-ci-robot](https://github.com/openshift-ci-robot) monitors
the queue and merges PRs that pass all of the required tests.

Only members of the
[Team OpenShift Ansible Committers](https://github.com/orgs/openshift/teams/team-openshift-ansible-committers)
can perform manual merges.

## Useful links

- Repository containing Jenkins job definitions: https://github.com/openshift/aos-cd-jobs
- List of required successful jobs before merge: https://github.com/openshift/aos-cd-jobs/blob/master/sjb/test_status_config.yml
- Source code of the bot responsible for testing and merging PRs: https://github.com/openshift/test-pull-requests/
- Trend of the time taken by merge jobs: https://ci.openshift.redhat.com/jenkins/job/merge_pull_request_openshift_ansible/buildTimeTrend
