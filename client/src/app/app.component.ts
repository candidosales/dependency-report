import { Component, OnInit } from '@angular/core';
import { DataService } from './providers/data.service';
import { GraphData } from './interface/data.interface';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

    graphData: GraphData = {
        projectsByFilters: {
            labels: [],
            datasets: {
                values: []
            }
        },
        componentsByFilters: {
            labels: [],
            datasets: {
                values: []
            }
        },
        componentsByProject: {
            labels: [],
            datasets: {
                values: []
            }
        },
        componentsByVersionAllProjects: {
            labels: [],
            datasets: {
                values: []
            }
        }
    }

    projects = [];
    components = [];

    componentsUsedByVersion = [
        {
            repository: {
                name: '@vendasta/uikit'
            },
            versions: [
                {
                    version: '8.0.1',
                    quantity: 10,
                    projects: ['listing-builder-client', 'customer-voice-client', 'concierge-cliet',
                    'reputation-client', 'iam-client', 'snapshot-client']
                },
                {
                    version: '9.0.1',
                    quantity: 2,
                    projects: ['snapshot-client', 'salesforce-client']
                },
                {
                    version: '7.0.1',
                    quantity: 1,
                    projects: ['vetl-client']
                }
            ]
        }, {
            repository: {
                name: '@vendasta/forms'
            },
            versions: [
                {
                    version: '8.0.1',
                    quantity: 10,
                    projects: ['listing-builder-client', 'customer-voice-client', 'concierge-cliet',
                    'reputation-client', 'iam-client', 'snapshot-client']
                },
                {
                    version: '9.0.1',
                    quantity: 2,
                    projects: ['snapshot-client', 'salesforce-client']
                },
                {
                    version: '7.0.1',
                    quantity: 1,
                    projects: ['vetl-client']
                }
            ]
        }
    ]

    constructor(private dataService: DataService) {}

    ngOnInit() {
        this.dataService.getDataInServer()
        .subscribe(value => {
            this.graphData = value.graphData;
            this.projects = value.projects;
            this.components = value.components;
        });
    }


    getProjectsTooltip(projects: Array<string>): string {
        return projects.join(' / ');
    }
}
