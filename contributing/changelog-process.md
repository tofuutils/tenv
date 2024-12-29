## Changelog Process

We use the [go-changelog](https://github.com/hashicorp/go-changelog) to generate and update the changelog from files created in the `.changelog/` directory. It is important that when you raise your Pull Request, there is a changelog entry which describes the changes your contribution makes. Not all changes require an entry in the CHANGELOG, guidance follows on what changes do.

### Changelog Format

The changelog format requires an entry in the following format, where HEADER corresponds to the changelog category, and the entry is the changelog entry itself. The entry should be included in a file in the `.changelog` directory with the naming convention `{PR-NUMBER}.txt`. For example, to create a changelog entry for pull request 1234, there should be a file named `.changelog/1234.txt`.

``````markdown
```release-note:{HEADER}
{ENTRY}
```
``````

If a pull request should contain multiple changelog entries, then multiple blocks can be added to the same changelog file. For example:

``````markdown
```release-note:note
The `broken` attribute has been deprecated. All configurations using `broken` should be updated to use the new `not_broken` attribute instead.
```

```release-note:enhancement
Add `not_broken` attribute
```
``````

### Skipping changelog entries

In order to skip/pass the automated checks where a CHANGELOG entry is not required, apply the `workflow/skip-changelog-entry` label.

### Pull Request Types to CHANGELOG

The CHANGELOG is intended to show operator-impacting changes to the codebase for a particular version. If every change or commit to the code resulted in an entry, the CHANGELOG would become less useful for operators. The lists below are general guidelines and examples for when a decision needs to be made to decide whether a change should have an entry.

#### Changes that should have a CHANGELOG entry

##### New full-length documentation guides (e.g. Getting Started Guide)

A new full length documentation entry gives the title of the documentation added, using the `release-note:new-guide` header.

``````markdown
```release-note:new-guide
How To Get Started With Tool X
```
``````

##### Version manager bug fixes

A new bug entry should use the `release-note:bug` header. 

``````markdown
```release-note:bug
Fix 'thing' being optional
```
``````

##### Version manager enhancements

A new enhancement entry should use the `release-note:enhancement` header.

``````markdown
```release-note:enhancement
Add new capability
```
``````

##### Deprecations

A deprecation entry should use the `release-note:note` header.

``````markdown
```release-note:note
X attribute is being deprecated in favor of the new Y attribute
```
``````

##### Breaking Changes and Removals

A breaking-change entry should use the `release-note:breaking-change` header.

``````markdown
```release-note:breaking-change
Resource no longer works for 'EXAMPLE' parameters
```
``````

#### Changes that may have a CHANGELOG entry

Dependency updates: If the update contains relevant bug fixes or enhancements that affect operators, those should be called out.
Any changes which do not fit into the above categories but warrant highlighting.

``````markdown
```release-note:note
Example resource now does X slightly differently
```

```release-note:dependency
`go-changelog` v0.1.0 => v0.1.1
```
``````

#### Changes that should _not_ have a CHANGELOG entry

- Testing updates
- Code refactoring (context dependent)
