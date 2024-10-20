import { Pipe, PipeTransform } from '@angular/core';
import memo from 'memo-decorator';

@Pipe({
    name: 'iconLibrary',
    standalone: true
})
export class IconLibraryPipe implements PipeTransform {

  @memo()
  transform(filter: string): string {
    return this.getIconLibraryByFilter(filter);
  }

  getIconLibraryByFilter(filter: string): string {
    if (filter.includes('angular')) {
      return 'https://coryrylan.com/assets/images/posts/types/angular.svg';
    }
    if (filter.includes('react')) {
      return 'https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg';
    }
    if (filter.includes('vue')) {
      return 'https://upload.wikimedia.org/wikipedia/commons/9/95/Vue.js_Logo_2.svg';
    }
    if (filter.includes('svelte')) {
      return 'https://upload.wikimedia.org/wikipedia/commons/1/1b/Svelte_Logo.svg';
    }
    return '';
  }
}
