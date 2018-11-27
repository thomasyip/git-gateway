# git-gateway - Gateway to hosted git APIs

**Secure role based access to the APIs of common Git Hosting providers.**

When building sites with a JAMstack approach, a common pattern is to store all content as structured data in a Git repository instead of relying on an external database.

Netlify CMS is an open-source content management UI that allows content editors to work with your content in Git through a familiar content editing interface. This allows people to write and edit content without having to write code or know anything about Git, markdown, YAML, JSON, etc.

However, for most use cases you won’t want to require all content editors to have an account with full access to the source code repository for your website.

Netlify’s Git Gateway lets you setup a gateway to your choice of Git provider's API ( now available with both GitHub and GitLab 🎉 ) that lets tools like Netlify CMS work with content, branches and pull requests on your users’ behalf.

The Git Gateway works with some supported identity service that can issue JWTs and only allows access when a JSON Web Token with sufficient permissions is present.

To configure the gateway, see our `example.env` file

The Gateway limits access to the following sub endpoints of the repository:

for GitHub:
```
   /repos/:owner/:name/git/
   /repos/:owner/:name/contents/
   /repos/:owner/:name/pulls/
   /repos/:owner/:name/branches/
```
for GitLab:
```
   /repos/:owner/:name/files/
   /repos/:owner/:name/commits/
   /repos/:owner/:name/tree/
```

**Running `git-gateway`**
**(Do not merge this section back to the open source project)**
**(Do not deploy it to production. It is a Proof of Concept has has not been secured. See @TODO items in code.**
**(the instruction assume Okta, and github.com)**

1. pull down this project
2. generate a `personal access token` on github. (recommended: using a test account and w/ `repo:status` and `public_repo` permission only)
    https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
3. `cp example.env my.env`
4. update `GITGATEWAY_GITHUB_ACCESS_TOKEN` value in `my.env` accordingly
5. update `GITGATEWAY_GITHUB_REPO` value in `my.env` (it will be where the content being stored, eg, `owner/netlify-test`.)
6. sign up for a Dev account on Okta: https://developer.okta.com/signup/
7. create a SPA Application onto the Dev account:
    a. fill out the details
    b. Pick "Send ID Token directly to app (Okta Simplified)"
    c. have redirect uri points to the url of your content-cms ip:port
      (eg, `http://localhost:8080/admin` etc, see, https://github.com/<< your org >>/content-cms)
8. update `ISSUER` value in `my.env` accordingly (eg, `https://dev-1234.oktapreview.com/oauth2/default`)
9. update `CLIENT_ID` value in `my.env` accordingly (eg, `32q897q234q324rq42322q`)
10. install Docker and add the `localdev` network
11. inspect Dockfile and then build the docker with this command:
    `docker build -t netlify/git-gateway:latest .`
12. run `git-gateway` with this command:
    `docker run --rm --env-file my.env --net localdev -p 127.0.0.1:8087:8087 --expose 8087 -ti --name netlify-git-gateway "netlify/git-gateway:latest"`
13. update `config.yml` in your content-cms repo (ie, https://github.com/<< your org >>/content-cms).
     change `backend.name` value to `git-gateway`
     change `backend.gateway_url` value to `http://localhost:8087`
14. run `content-cms` following the README.md

**Develop, Build and Run git-gateway**

1. Follow instructions 1 - 10 in previous "Running `git-gateway`" section
2. Run these commands once:
   ```
   docker build -t netlify/git-gateway:latest .
   docker run --rm --env-file my.env --net localdev -p 127.0.0.1:8087:8087 --expose 8087 -ti -v $PWD:/go/src/github.com/netlify/git-gateway --entrypoint '/bin/sh' --user root netlify/git-gateway:latest
   cd /go/src/github.com/netlify/git-gateway
   make deps
   ```
3. Run these commands after edit:
   ```
   make build && ./git-gateway
   ```
4. `<ctrl> + c` to stop

