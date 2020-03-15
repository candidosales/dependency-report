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

    treemap = [
        ['Location', 'Parent', 'Market trade volume (size)'],
        ['@angular/core',    null,                 0],
        ['@angular/core_9',    '@angular/core',             0],
        ['@angular/core_8',    '@angular/core',             0],
        ['@angular/core_7',    '@angular/core',             0],
        ['Brazil',    '@angular/core_9',            11],
        ['USA',       '@angular/core_9',            52],
        ['Mexico',    '@angular/core_9',            24],
        ['Canada',    '@angular/core_9',            16],
        ['France',    '@angular/core_8',             42],
        ['Germany',   '@angular/core_8',             31],
        ['Sweden',    '@angular/core_8',             22],
        ['Italy',     '@angular/core_8',             17],
        ['UK',        '@angular/core_8',             21],
        ['China',     '@angular/core_7',               36],
        ['Japan',     '@angular/core_7',               20],
        ['India',     '@angular/core_7',               40],
        ['Laos',      '@angular/core_7',               4],
        ['Mongolia',  '@angular/core_7',               1],
        ['Israel',    '@angular/core_7',               12],
        ['Iran',      '@angular/core_7',               18],
        ['Pakistan',  '@angular/core_7',               11],
      ];

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
