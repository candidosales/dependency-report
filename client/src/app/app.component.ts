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
    components = []

    constructor(private dataService: DataService) {}

    ngOnInit() {
        this.dataService.getDataInServer()
        .subscribe(value => {
            this.graphData = value.graphData;
            this.projects = value.projects;
            this.components = value.components;
        });
    }
}
