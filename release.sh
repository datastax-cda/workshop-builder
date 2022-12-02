RELVER=1.5

# Tag release with an appropriate version and comment
echo Before uploading any release artifacts, tag the release with an appropriate version and comment as shown, then update the version in this script and run.
echo git tag -a v${RELVER} -m \"Release comments...\"
echo git push origin v${RELVER}
echo

# Fetch deps and build Mac release
echo Building Mac binary...
GOOS=darwin GOARCH=amd64 go get ./...
GOOS=darwin GOARCH=amd64 go build -o release/pace-darwin

# Fetch deps and build Linux release
echo Building Linux binary...
GOOS=linux GOARCH=amd64 go get ./...
GOOS=linux GOARCH=amd64 go build -o release/pace-linux 

# Fetch deps and build Windows release
echo Building Windows binary...
GOOS=windows GOARCH=amd64 go get ./...
GOOS=windows GOARCH=amd64 go build -o release/pace-windows.exe 

# Create source archives
# Archives are created automatically by github when a new release is published,
# so we no longer need to create them manually
#echo Creating zip source archive...
#zip -r release/pace-builder-${RELVER}.zip . -x .git/\* .gitignore release/\*
#echo Creating tar source archive...
#tar -cvzf release/pace-builder-${RELVER}.tar.gz --exclude .git/\* --exclude .gitignore --exclude release/\* .

echo -e "\nArtifacts in release directory can now be uploaded to https://github.com/Pivotal-Field-Engineering/pace-builder/releases"

echo -e "\nOptionally update the brew tap.  Create a fork of the repo below, then update and commit the specified file with the noted changes and create a pull request on the original repo using compare across forks."
echo "Changes to be made:"
echo "  URL: https://github.com/pivotal-legacy/homebrew-tap"
echo "  File: pace-cli.rb"
echo "  Version: ${RELVER}"
echo "  SHA256: `openssl dgst -sha256 release/pace-darwin | cut -d\  -f2`"
