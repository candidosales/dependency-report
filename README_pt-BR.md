# 📊 Relatório de dependências

Relatório de dependências visa ajudar a analisar a consistência das dependências em seus projetos frontend de sua empresa. Veja um exemplo [https://dependency-report.web.app/](https://dependency-report.web.app/). [In English](./README.md)

## 🎯 Recursos

- Visualizar qual a porcentagem de projetos estão em uma determinada versão;
- Visualizar qual a porcentagem de componentes estão em uma determinada versão;
- Visualizar todas as dependências de seus projetos e quantas diferentes versões são usadas por projetos;
- Visualizar link para documentacão de frontend e design para cada projeto/componente;
- Visualizar notificacões do Github relacionados a vulnerabilidades de seguranca para cada projeto/componente;
- Visualizar quais dependências precisam ser atualizadas por projetos;

![Github Notifications](https://media.giphy.com/media/kfFGCtQ8m1M8hF8qKN/giphy.gif)

## ⚙️ Como configurar

O projeto é dividido em duas pastas:

- `client`: Frontend desenvolvido em Angular para apresentacão dos gráficos. Angular 12+ - [Como instalar Node?](https://nodejs.org/en/download/package-manager/);
- `server`: Backend desenvolvido em Go para geracão dos dados para os relatórios. Go 1.15 - [Como instalar Go?](https://golang.org/doc/install);

### Criando seu personal Github Auth Token

É necessário criar seu `Personal Access Token` para que o servidor tenha permissão de usar a API do Github. Para mais instrucões [acesse](https://docs.cachethq.io/docs/github-oauth-token) ou [aqui](https://github.com/settings/tokens).

### Inicializando o cliente

O projeto usa Angular 9+ e requer versão do Node 12+. Você pode usar o [NVM](https://github.com/nvm-sh/nvm) para controlar as versões node de sua máquina.

```bash
cd client/
npm install # Instalar dependências
ng serve
```

Irá inicializar na porta `4200`. [http://localhost:4200](http://localhost:4200).

### Inicializando o servidor

```bash
cd server/
GITHUB_AUTH_TOKEN=<personal-auth-token> go run *.go # Exemplo: GITHUB_AUTH_TOKEN=12321wqdd12e12321dse go run *.go
```

Irá criar um servidor na porta `3000`. [http://localhost:3000](http://localhost:3000).

### Adaptando para seus projetos

Você precisa editar o arquivo [server/config.json](./server/config.json) adicionar os repositórios de seus projetos frontend e components. Não esqueca de especificar o tipo de repositório se é `project` ou `component`.

- `project`: São os seus projetos frontend que podem ser um admin ou backoffice de sua empresa;
- `component`: São os componentes modulares que são usados em seus projetos;

Exemplo:

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

O filtro é usado para analisar as dependências de seus repositórios e assim gerar as estatísticas. O formato do filtro é o nome da biblioteca e a versão: `<nome-biblioteca>_<versão>`. Você ver sua dependência no seu `package.json` nesse formato `"@angular/core": "~9.1.1"` então você cria seu filtro dessa forma: `@angular/core_9.1.1`.

Dessa forma, você criar qualquer filtro para qualquer dependência ou versão.

Exemplo:

- `"@angular/core": "~9.0.0"` mude para `@angular/core_9.0.0` ou `@angular/core_9`;
- `"react": "^16.12.0"` mude para `react_16.12.0` ou `react_16`;

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

### Inicializando

Após a configuracão do [server/config.json](./server/config.json) e inicializacão do `server` e `client` você acessa o `client` via [http://localhost:4200](http://localhost:4200) e clica no botão de `REFRESH` na barra do topo para gerar os dados que você precisa.

## 📚 Referência

- [https://segment.com/blog/driving-adoption-of-a-design-system/](https://segment.com/blog/driving-adoption-of-a-design-system/)

## 👍 Contribuir

Se você quer agradecer e/ou apoiar o desenvolvimento deste projeto:

1. Adicione uma [GitHub Star](https://github.com/candidosales/dependency-report/stargazers) para o projeto.
2. Tweet sobre o projeto [em seu Twitter](https://twitter.com/intent/tweet?url=https%3A%2F%2Fgithub.com%2Fcandidosales%2Fdependency-report&text=Dependencies%20report%20aims%20to%20help%20analyze%20the%20consistency%20of%20the%20dependencies%20in%20your%20company%27s%20frontend%20projects).
3. Escreva um review ou tutorial no [Medium](https://medium.com/), [Dev.to](https://dev.to/) ou no blog pessoal.
4. Apoie o projeto doando um [cafézinho](https://buymeacoff.ee/candidosales).

## ☕ Apoiadores

Se você quer apoiar o Relatório de dependência, você pode ☕ [**comprar um café aqui**](https://buymeacoff.ee/candidosales)

## ✨ Agradecimentos

[Joel Kesler](https://github.com/joelkesler), [Rafael Lins](https://github.com/g0dkar)

## Autor

- Cândido Sales - [@candidosales](https://twitter.com/candidosales)

## ⚠️ Direitos autorais e licenca

Direitos autorais de código e documentação 2020-2030, os [Autores](https://github.com/candidosales/dependency-report/graphs/contributors) e o Código liberado sob a [MIT License](https://github.com/candidosales/dependency-report/blob/master/LICENSE). Documentos publicados sob [Creative Commons](https://creativecommons.org/licenses/by/3.0/).
