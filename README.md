# google-photos-cleaner

The idea is to clean all memes and screenshots from the google photos to save some space.

There are multiple catches on how google has exposed its photos API endpoints.
  - There is no delete endpoint exposed (There is an API to delete files from an album if the app thats deleting it has created the files/album)
  - There is no move endpoint exposed. So we cant move all the photos to a new album and manually delete it

To run this code, you would need a credentials.json file which you can get from your google developers console page

**This is not tested on Windows**


test
