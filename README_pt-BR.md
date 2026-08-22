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

- `web`: Frontend desenvolvido em React (TanStack Start) para apresentacão dos gráficos. Node 20+ - [Como instalar Node?](https://nodejs.org/en/download/package-manager/).
- `server`: Backend desenvolvido em Go para geracão dos dados para os relatórios. Go 1.27+ - [Como instalar Go?](https://golang.org/doc/install);

### Criando seu personal Github Auth Token

É necessário criar seu `Personal Access Token` para que o servidor tenha permissão de usar a API do Github. Para mais instrucões [acesse](https://docs.cachethq.io/docs/github-oauth-token) ou [aqui](https://github.com/settings/tokens).

### Inicializando o web app

```bash
cd web/
pnpm install # Instalar dependências
pnpm dev
```

Irá inicializar na porta `3000`. [http://localhost:3000](http://localhost:3000).

O web app lê a URL do servidor em `web/.env` (`VITE_API_URL`, padrão `http://127.0.0.1:8080`). Ajuste caso seu servidor rode em outro host/porta.

### Inicializando o servidor

```bash
cd server/
GITHUB_AUTH_TOKEN=<personal-auth-token> go run . # Exemplo: GITHUB_AUTH_TOKEN=12321wqdd12e12321dse go run .
```

Irá criar um servidor na porta `8080`. [http://localhost:8080](http://localhost:8080).

### Adaptando para seus projetos

Você precisa editar o arquivo [server/config.json](./server/config.json) adicionar os repositórios de seus projetos frontend e components. Não esqueca de especificar o tipo de repositório se é `project` ou `component`.

- `project`: São os seus projetos frontend que podem ser um admin ou backoffice de sua empresa;
- `component`: São os componentes modulares que são usados em seus projetos;

Exemplo:

```json
{
    "repositories": [
        {
            "url": "https://github.com/candidosales/guide-mental-health",
            "type": "project"
        },
        {
            "url": "https://github.com/candidosales/dependency-report/tree/main/web",
            "type": "component"
        },
    ]
}
```

O filtro é usado para analisar as dependências de seus repositórios e assim gerar as estatísticas. O formato do filtro é o nome da biblioteca e a versão: `<nome-biblioteca>_<versão>`. Você ver sua dependência no seu `package.json` nesse formato `"react":"~19.1.1"` então você cria seu filtro dessa forma: `react_19.1.1`.

Dessa forma, você criar qualquer filtro para qualquer dependência ou versão.

Exemplo:

- `"react": "^19.0.0"` mude para `react_19.0.0` ou `react_19`;
- `"next": "^16.12.0"` mude para `next_16.12.0` ou `next_16`;

```json
{
    "filters": [
        "react_16",
        "react_17",
        "react_18",
        "react_19"
    ],
}
```

### Inicializando

Após a configuracão do [server/config.json](./server/config.json) e inicializacão do `server` e `web` você acessa o `web` via [http://localhost:3000](http://localhost:3000) e clica no botão de `REFRESH` na barra do topo para gerar os dados que você precisa.

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

Direitos autorais de código e documentação 2024-2034, os [Autores](https://github.com/candidosales/dependency-report/graphs/contributors) e o Código liberado sob a [MIT License](https://github.com/candidosales/dependency-report/blob/master/LICENSE). Documentos publicados sob [Creative Commons](https://creativecommons.org/licenses/by/3.0/).
