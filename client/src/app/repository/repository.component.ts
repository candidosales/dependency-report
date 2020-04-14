import { Component, ChangeDetectionStrategy, Input } from '@angular/core';
import { Repository } from '../interface/repository.interface';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RepositoryComponent {
    @Input() repository: Repository;
}