# 📊 Dependencies report

Dependencies report aims to help analyze the consistency of the dependencies in your company's frontend projects. See an example [https://dependency-report.web.app/](https://dependency-report.web.app/). [Em português](./README_pt-BR.md)

## 🎯 Features

- View what percentage of projects are in a given version;
- View what percentage of components are in a given version;
- View all the dependencies of your projects and how many different versions are used by projects;
- View link to frontend and design documentation for each project/component;
- View Github notifications related to security vulnerabilities for each project/component;
- View which dependencies need to be updated by projects;

![Github Notifications](https://media.giphy.com/media/kfFGCtQ8m1M8hF8qKN/giphy.gif)

## ⚙️ How to setup

The project is divided into two folders:

- `client`: Frontend developed in Angular for presentation of graphics. Angular 9+ - [How to install Node?](https://nodejs.org/en/download/package-manager/);
- `server`: Backend developed in Go to generate data for reports. Go 1.14 - [How to install Go?](https://golang.org/doc/install);

### Creating your personal Github Auth Token

It is necessary to create your `Personal Access Token` for the server to be allowed to use the Github API. For further instructions [visit](https://docs.cachethq.io/docs/github-oauth-token) or [here](https://github.com/settings/tokens).

### Initializing the client

The project uses Angular 9+ and requires version of Node 12+. You can use the [NVM](https://github.com/nvm-sh/nvm) to control the node versions of your machine.

```bash
cd client/
npm install # Install dependencies
ng serve
```

Will boot on the port `4200`. [http://localhost:4200](http://localhost:4200).

### Initializing the server

```bash
cd server/
GITHUB_AUTH_TOKEN=<personal-auth-token> go run *.go # Example: GITHUB_AUTH_TOKEN=12321wqdd12e12321dse go run *.go
```

Will create a server on the port `3000`. [http://localhost:3000](http://localhost:3000).

### Adapting to your projects

You have to edit the [server/config.json](./server/config.json) file to add the repositories for your frontend and component projects. Don't forget to specify the type of repository if it is `project` or `component`.

- `project`: Are your frontend projects that can be an admin or backoffice for your company;
- `component`: Modular components that are used in your projects;

Example:

```json
{
    "repositories": [
        {
            "url": "https://github.com/vendasta/listing-builder-client",
            "type": "project"
        },
        {
            "url": "https://github.com/vendasta/frontend/tree/master/angular/projects/business-categories",
            "type": "component"
        },
    ]
}
```

The filter is used to analyze the dependencies of your repositories and thus generate the statistics. The filter format is the name of the library and the version: `<library-name>_<version>`. You see your dependency on your `package.json` in this format `"@angular/core":"~9.1.1"` then you create your filter this way: `@angular/core_9.1.1`.

That way, you can create any filter for any dependency or version..

Example:

- `"@angular/core": "~9.0.0"` switch to `@angular/core_9.0.0` or `@angular/core_9`;
- `"react": "^16.12.0"` switch to `react_16.12.0` or `react_16`;

```json
{
    "filters": [
        "@angular/core_4",
        "@angular/core_6",
        "@angular/core_7",
        "@angular/core_8",
        "@angular/core_9"
    ],
}
```

### Initializing

After configuring [server/config.json](./server/config.json) and initializing `server` and `client` you access `client` via [http://localhost:4200](http://localhost:4200) and click the button `REFRESH` on the top bar to generate the data you need.

## 📚 Reference

- [https://segment.com/blog/driving-adoption-of-a-design-system/](https://segment.com/blog/driving-adoption-of-a-design-system/)

## 👍 Contribute

If you want to say thank you and/or support the active development this project:

1. Add a [GitHub Star](https://github.com/candidosales/dependency-report/stargazers) to the project.
2. Tweet about the project [on your Twitter](https://twitter.com/intent/tweet?url=https%3A%2F%2Fgithub.com%2Fcandidosales%2Fdependency-report&text=Dependencies%20report%20aims%20to%20help%20analyze%20the%20consistency%20of%20the%20dependencies%20in%20your%20company%27s%20frontend%20projects).
3. Write a review or tutorial on [Medium](https://medium.com/), [Dev.to](https://dev.to/) or personal blog.
4. Support the project by donating a [cup of coffee](https://buymeacoff.ee/candidosales).

## ☕ Supporters

If you want to support Dependency Report, you can ☕ [**buy a coffee here**](https://buymeacoff.ee/candidosales)

## ✨ Acknowledgment

[Joel Kesler](https://github.com/joelkesler), [Rafael Lins](https://github.com/g0dkar)

## Author

- Cândido Sales - [@candidosales](https://twitter.com/candidosales)

## ⚠️ Copyright and license

Code and documentation copyright 2020-2030 the [Authors](https://github.com/candidosales/dependency-report/graphs/contributors) and Code released under the [MIT License](https://github.com/candidosales/dependency-report/blob/master/LICENSE). Docs released under [Creative Commons](https://creativecommons.org/licenses/by/3.0/).