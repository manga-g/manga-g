#!/usr/bin/env bash

# Removes the old builds within the project
echo "Removing old builds"
rm ../cmd/core/MangaG
rm ../cmd/core/manga-g

# Removes if there is a binary file from bin directory
echo "Removing old manga-g version file"
rm /usr/local/bin/manga-g

# Builds the project and creates a binary file
echo "Building new version of manga-g"
# let the install script run the build script
chmod +x build.sh
sh build.sh

# Copy the binary file to the usr bin directory
echo "Installing manga-g"
cp ../manga-g /usr/local/bin/manga-g
echo "Install has finished"