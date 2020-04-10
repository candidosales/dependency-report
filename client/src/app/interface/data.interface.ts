import { Repository } from './repository.interface';
export interface Data  {
    generatedAt: string;
    projects: Array<Repository>;
    components: Array<Repository>;
    dependenciesByVersions: any;
    graphData: GraphData;
}

export interface GraphData {
    componentsByFilters: [[string, any]];
    projectsByFilters: [[string, any]];

}

