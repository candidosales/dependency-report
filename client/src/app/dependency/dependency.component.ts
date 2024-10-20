import { Component, ChangeDetectionStrategy, input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { MatTooltip } from '@angular/material/tooltip';

@Component({
    selector: 'app-dependency',
    templateUrl: './dependency.component.html',
    styleUrls: ['./dependency.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true,
    imports: [MatTooltip, TranslateModule]
})
export class DependencyComponent {

    name = input<string>();
    dependency = input<any>();

    objectKeys = Object.keys;

    getCountObjectKeys(object: any): number {
        return Object.keys(object).length;
    }
}