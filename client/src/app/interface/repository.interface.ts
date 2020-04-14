export interface Repository {
    name: string;
    version: string;
    filter: string;
    url: string;
    documentation: Documentation;
}

export interface Documentation {
    frontend: string;
    design: string;
}