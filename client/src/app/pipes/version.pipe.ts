import { Pipe, PipeTransform } from '@angular/core';
import memo from 'memo-decorator';

@Pipe({
    name: 'version',
    standalone: true
})
export class VersionPipe implements PipeTransform {

  @memo()
  transform(filter: string): string {
    return this.getVersionByFilter(filter);
  }

  getVersionByFilter(filter: string): string {
    const values = filter.split('_');
    return values[1];
  }
}
