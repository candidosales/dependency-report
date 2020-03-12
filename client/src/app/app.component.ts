import { Component, OnInit } from '@angular/core';
import { DataService } from './providers/data.service';
import { GraphData } from './interface/data.interface';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

    graphDataGoogle: GraphData;

    projects = [];
    components = [];

    objectKeys = Object.keys;

    componentsByVersions = {};

    constructor(private dataService: DataService) {}

    ngOnInit() {
        this.dataService.getDataInServer()
        .subscribe(value => {
            value.graphData.projectsByFilters?.unshift(['Filter', 'Version']);
            value.graphData.componentsByFilters?.unshift(['Filter', 'Version']);

            this.graphDataGoogle = value.graphData;

            this.projects = value.projects;
            this.components = value.components;

            this.componentsByVersions = value.componentsByVersions;
        });
    }


    getProjectsTooltip(projects: Array<string>): string {
        return projects.join(' / ');
    }

    getVersionByFilter(filter: string): string {
        const values = filter.split('_');
        return values[1];
    }

    getIconByFilter(filter: string): string {
        if (filter.includes('angular')) {
            return 'https://coryrylan.com/assets/images/posts/types/angular.svg';
        }
        return '';
    }
}
