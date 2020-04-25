import { Repository } from './repository.interface';
export interface Data  {
    generatedAt: string;
    summary: SummaryData;
    projects: Array<Repository>;
    components: Array<Repository>;
    dependenciesByVersions: any;
    graphData?: GraphData;
}

export interface SummaryData {
    updated: Array<string>;
    inconsistent: Array<string>;
    vulnerable: Array<string>;
}
export interface GraphData {
    componentsByFilters: [[string, any]];
    projectsByFilters: [[string, any]];
}

