import { Component, ChangeDetectionStrategy, Input } from '@angular/core';
import { Repository } from '../interface/repository.interface';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RepositoryComponent {
    @Input() repository: Repository;

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