import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-dependency',
  templateUrl: './dependency.component.html',
  styleUrls: ['./dependency.component.scss']
})
export class DependencyComponent {

    @Input() name: string;
    @Input() dependency: any;

    objectKeys = Object.keys;

    getProjectsTooltip(projects: Array<string>): string {
        return projects.join(' / ');
    }

    getCountObjectKeys(object: any): number {
        return Object.keys(object).length;
    }

    getClassByDependenciesVersions(quantity: number): string {
        if (quantity <= 2) {
            return 'repository__versions--good';
        }
        if (quantity > 2 && quantity <= 5) {
            return 'repository__versions--warning';
        }
        if (quantity > 5 && quantity <= 10) {
            return 'repository__versions--bad';
        }
        if (quantity > 10) {
            return 'repository__versions--terrible';
        }
        return '';
    }
}