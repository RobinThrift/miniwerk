# generate a release with the given tag
release tag:
    just changelog {{tag}}
    git add CHANGELOG
    git commit -m "Releasing version {{tag}}"
    git tag {{tag}}
    git push
    git push origin {{tag}}

# generate a changelog using github.com/git-chglog/git-chglog
changelog tag: 
    git-chglog -o CHANGELOG/CHANGELOG-{{tag}}.md --next-tag {{tag}} --sort semver --config CHANGELOG/config.yml {{ tag }}
    echo "- [CHANGELOG-{{tag}}.md](./CHANGELOG-{{tag}}.md)" >> CHANGELOG/README.md

