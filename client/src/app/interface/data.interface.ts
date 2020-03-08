import { Repository } from './repository.interface';
export interface Data  {
    projects: Array<Repository>;
    components: Array<Repository>;
    graphData: GraphData;
}

export interface GraphData {
    componentsByProject: StatsDataFrappe;
    componentsByVersionAllProjects: StatsDataFrappe;
    projectsByFilters: StatsDataFrappe;
    componentsByFilters: StatsDataFrappe;
}

export interface StatsDataFrappe {
    labels: Array<string>;
    datasets: StatsDataset;
}

export interface StatsDataset {
    values: Array<number>;
}

