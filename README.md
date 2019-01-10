# medat-wortfluessigkeit-backend

This is the backend for [https://medat-wortfluessigkeit.appspot.com](). The frontend can be found in this repository: [https://github.com/schoenenberg/medat-wortfluessigkeit-frontend]()

To publish the complete app on Google App Engine:
* Run `npm run build` in frontend-repository
* copy contents of `dist/` directory to this repositories `frontend/` directory
* finally run `gcloud app deploy`

Do not forget to create a storage bucket an upload the `words.csv`. 
