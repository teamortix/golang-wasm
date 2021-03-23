# Contributing to Golang-WASM
If you're reading this, thanks for taking a look! We have a couple of guidelines that we encourage you to follow when working on Golang-WASM.

## Code of Conduct
If you find any bugs or would like to request any features, please open an issue. Once the issue is open, feel free to put in a pull request to fix the issue. Small corrections, like fixing spelling errors, or updating broken urls are appreciated. 

We encourage an open, friendly, and supportive environment around the development of Golang-WASM. If you disagree with someone for any reason, discuss the issue and express you opinions, don't attack the person. Discrimination of any kind against any person is not permitted. If you detract from this project's collaborative environment, you'll be prevented from participating in the future development of this project until you prove you can behave yourself adequately.

> Please use sound reasoning to support your suggestions - don't rely on arguments based on 'years of experience,' supposed skill, job title, etc. to get your points across.

# General Guidelines
All exported functions, variables, and constants must contain documentation. Readable code with clear behavior works better than illegible optimized code. Use comments for unexported functions when their purpose is non-trivial.

Each commit should denote exactly one thing, whether it be a bug fix, or a feature addition. Try not to do both at the same time - it makes code harder to review. Once the codebase is stable, the contributing flow should be:

1. Check for issues and discussion on a bug/feature.
2. If none found, open an issue about.
3. Fork and develop the feature.
4. Open a pull request. If an issue is related, tag the issue number in the commit message
5. Merge the code into the project.

Once the API for Golang-WASM reaches a stable level, SemVer will be used for tagging new bug fixes, features, and breaking changes. 

> Our project uses the Conventional Commits standard for writing commit messages. Commits that do not follow this will not be merged into the code base. Read more about Conventional Commits [here](https://conventionalcommits.org)