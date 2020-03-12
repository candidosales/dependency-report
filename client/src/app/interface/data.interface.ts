import { Repository } from './repository.interface';
export interface Data  {
    projects: Array<Repository>;
    components: Array<Repository>;
    componentsByVersions: any;
    graphData: GraphData;

}

export interface GraphData {
    componentsByFilters: [[string, any]];
    projectsByFilters: [[string, any]];

}

