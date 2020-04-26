import { Component, Input, ChangeDetectionStrategy } from '@angular/core';

@Component({
  selector: 'app-dependency',
  templateUrl: './dependency.component.html',
  styleUrls: ['./dependency.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
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
}